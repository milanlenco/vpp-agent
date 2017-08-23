#!/usr/bin/python

import sys, getopt, os
import yaml
import subprocess


AGENT_PCKG = "github.com/ligato/vpp-agent"
CN_INFRA_PCKG = "github.com/ligato/cn-infra"
CACHE = ".resolve-deps.cache"
CACHE_CONFLICT_RES = "conflict-resolution"


def get_git_revision_hash(path):
    return subprocess.check_output(['git', '-C', path, 'rev-parse', 'HEAD']).strip()


def load_yaml(path):
    with open(path, 'r') as ymlfile:
        data = yaml.load(ymlfile)
        return data


def dump_yaml(path, data):
    ymlfile = file(path, 'w')
    yaml.dump(data, ymlfile, default_flow_style=False)


def find_dependency(deps, package):
    target_dep = None
    for dep in deps:
        if dep.has_key('package') and dep['package'] == package:
            target_dep = dep
            break
    return target_dep


def add_dependency(target_deps, dep, quiet=False):
    if not quiet:
        package = dep['package']
        version = dep['version'] if dep.has_key('version') else 'unspecified'
        print "Inserting new dependency: '%s', version: '%s'" % (package, version)
    target_deps.append(dep)


def update_dependency(dep, new_dep, quiet=False):
    package = dep['package']
    if package != new_dep['package']:
        print "WARNING: dependency mismatch!"
    if not quiet:
        old_version = dep['version'] if dep.has_key('version') else 'unspecified'
        new_version = new_dep['version'] if new_dep.has_key('version') else 'unspecified'
        print "Updating dependency: '%s', version: '%s', to version: '%s'" % (package, old_version, new_version)
    dep.clear()
    dep.update(new_dep)


def prioritize_aliased_deps(agent_deps, target_deps):
    for dep in agent_deps:
        if dep.has_key('package') and dep.has_key('repo'):
            dest_dep = find_dependency(target_deps, dep['package'])
            if dest_dep is not None:
                update_dependency(dest_dep, dep)


def resolve_conflict(target_pckg, target_dep, agent_dep, conflicts_cache, use_cache):
    cached_deps = conflicts_cache['dependencies']
    cached_dep = find_dependency(cached_deps, target_dep['package'])
    if use_cache and cached_dep is not None:
        update_dependency(target_dep, cached_dep)
    else:
        prioritize_agent = False
        print "Do you want to prioritize agent over %s? (Y/N):" % target_pckg
        while True:
            choice = raw_input().lower()
            if choice == 'y':
                prioritize_agent = True  
                break
            elif choice == 'n':
                prioritize_agent = False
                break
            else:
                print "Please respond with 'Y/y' or 'N/n'"
        to_cache = target_dep
        if prioritize_agent:
            update_dependency(target_dep, agent_dep)
            to_cache = agent_dep    
        if cached_dep is None:
            add_dependency(cached_deps, to_cache, True)
        else:
            update_dependency(cached_dep, to_cache, True)


def resolve_conflicts(target_pckg, target_deps, agent_deps, conflicts_cache, use_cache):
    for dep in target_deps:
        dep_ver = dep['version'] if dep.has_key('version') else ''
        if dep.has_key('package'):
            package = dep['package']
            agent_dep = find_dependency(agent_deps, package)
            if agent_dep is not None:
                agent_dep_ver = agent_dep['version'] if agent_dep.has_key('version') else ''
                if dep_ver != agent_dep_ver:
                    print ("Conflicting dependency: '%s', target version: '%s', version used by agent: '%s'"
                            % (package, dep_ver, agent_dep_ver))
                    resolve_conflict(target_pckg, dep, agent_dep, conflicts_cache, use_cache)


def resolve_deps(agent, target, use_cache):
    script_path = os.path.dirname(os.path.realpath(__file__))
    target_pckg = "/".join(target.split("/")[-3:])

    # load glide.yaml file from the target, vpp-agent and cn-infra
    agent_glide_yaml = load_yaml(os.path.join(agent, "glide.yaml"))
    if not agent_glide_yaml.has_key('import'):
        agent_glide_yaml['import'] = []
    agent_deps = agent_glide_yaml['import']

    cn_infra_glide_yaml = load_yaml(os.path.join(agent, "vendor", CN_INFRA_PCKG, "glide.yaml"))
    if not cn_infra_glide_yaml.has_key('import'):
        cn_infra_glide_yaml['import'] = []
    cn_infra_deps = cn_infra_glide_yaml['import']
    
    target_glide_yaml = load_yaml(os.path.join(target, "glide.yaml"))
    if not target_glide_yaml.has_key('import'):
        target_glide_yaml['import'] = []
    target_deps = target_glide_yaml['import']

    # load or initialize the cache
    cache_yaml = {}
    try:
        cache_yaml = load_yaml(os.path.join(script_path, CACHE))
    except IOError:
        pass
    if not cache_yaml.has_key(CACHE_CONFLICT_RES):
        cache_yaml[CACHE_CONFLICT_RES] = []
    conflicts_cache = None
    for pckg in cache_yaml[CACHE_CONFLICT_RES]:
        if pckg.has_key('package') and pckg['package'] == target_pckg:
            conflicts_cache = pckg
            if not conflicts_cache.has_key("dependencies"):
                conflicts_cache['dependencies'] = []
            break
    if conflicts_cache is None:
        conflicts_cache = {'package': target_pckg, 'dependencies': []}
        cache_yaml[CACHE_CONFLICT_RES].append(conflicts_cache)

    print "Cached conflict resolutions: %s, use cache: %s" % (str(conflicts_cache), str(use_cache))
    #print "Agent revision: %s" % get_git_revision_hash(agent)
    #print "Agent deps: %s" % str(agent_deps)
    #print "Target glide yaml: %s" % str(target_glide_yaml)

    # explicitly specify version of the agent
    agent_dep = find_dependency(target_deps, AGENT_PCKG)
    desired_agent_dep = {'version': get_git_revision_hash(agent), 'package': AGENT_PCKG}
    if agent_dep is None:
        add_dependency(target_deps, desired_agent_dep)
    else:
        update_dependency(agent_dep, desired_agent_dep)

    # take cn-infra from the agent
    cn_infra_dep = find_dependency(target_deps, CN_INFRA_PCKG)
    desired_cn_infra_dep = find_dependency(agent_deps, CN_INFRA_PCKG)
    if cn_infra_dep is None:
        add_dependency(target_deps, desired_cn_infra_dep)
    else:
        update_dependency(cn_infra_dep, desired_cn_infra_dep)

    # prioritize *aliased* dependencies from agent and cn-infra
    prioritize_aliased_deps(agent_deps, target_deps)
    prioritize_aliased_deps(cn_infra_deps, target_deps)

    # resolve conflicts
    resolve_conflicts(target_pckg, target_deps, agent_deps, conflicts_cache, use_cache)
    resolve_conflicts(target_pckg, target_deps, cn_infra_deps, conflicts_cache, use_cache)

    print "Target glide yaml: %s" % str(target_glide_yaml)
    print
    print "Cache yaml: %s" % str(cache_yaml)

    # dump resolved version of target glide.yaml and the cache
    dump_yaml(os.path.join(target, "glide-resolved.yaml"), target_glide_yaml)
    dump_yaml(os.path.join(script_path, CACHE), cache_yaml)



def print_usage():
    print 'resolve-deps.py -a/--agent <agent-path> -t/--target <target-path> -n/--no-cache'


def main(argv):
    agent = ''
    target = ''
    use_cache = True
    try:
        opts, args = getopt.getopt(argv,"ha:t:n",["agent=","target=","no-cache"])
    except getopt.GetoptError:
        print_usage()
        sys.exit(2)
    for opt, arg in opts:
        if opt == '-h':
            print_usage
            sys.exit()
        elif opt in ("-a", "--agent"):
            agent = arg
        elif opt in ("-t", "--target"):
            target = arg
        elif opt in ("-n", "--no-cache"):
            use_cache = False
    if len(agent) == 0 or len(target) == 0:
        print_usage()
        sys.exit(2)
    resolve_deps(agent, target, use_cache)


if __name__ == "__main__":
   main(sys.argv[1:])
