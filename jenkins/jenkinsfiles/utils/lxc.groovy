/* groovylint-disable CompileStatic, LineLength */
void launchContainer(String release) {
    // sudo -u ${env.LXC_USER} -H -E -- 
    sh """
        lxc launch --ephemeral ${release} ${env.BUILD_TAG}
    """
    s('while [ ! -f /var/lib/cloud/instance/boot-finished ]; do sleep 0.1; done')
}

void pushWorkspace() {
    echo 'pushing workspace'
    // sudo -u $LXC_USER -H -E -- 
    sh '''
        mkdir tarball && tar --exclude=./tarball -czf  tarball/workspace.tar.gz .
        lxc file push  ./tarball/workspace.tar.gz $BUILD_TAG/home/ubuntu/${JOB_BASE_NAME}/ --create-dirs
        rm -rf tarball
    '''
    echo 'untarring'
    s('tar xzvf workspace.tar.gz')
    echo 'done'
    // s('ls -lah')
    // Can't use --uid/--gid/--mode in lxc file push recursive mode
    // So we just chown it after the fact.
    // s('sudo chown -R ubuntu:root ./')
}
//
void pullFileFromHome(String path) {
    // sudo -u ${env.LXC_USER} -H -E -- 
    cmd =  "lxc file pull ${env.BUILD_TAG}/home/ubuntu/${path}" 
    sh "${cmd}"
}

void s(String command, List<String> envArgs=[]) {
    // sudo -u ${env.LXC_USER} -H -E -- 
    // --env HTTP_PROXY==${env.HTTP_PROXY} --env HTTPS_PROXY=${env.HTTPS_PROXY}
    String cmd = "lxc exec ${env.BUILD_TAG} --cwd /home/ubuntu/${env.JOB_BASE_NAME} --user 1000 --env HOME=/home/ubuntu "
    envArgs.each { env -> cmd <<= (' --env ' + env + ' ') }
    cmd <<= ' -- bash -c '
    cmd <<= "\'${command.trim()}\'"
    echo "${cmd}"
    sh(script: "${cmd}")
}

void installSnap(String name, Boolean classic=false, String channel='') {
    String cmd = 'sudo snap install '
    cmd <<= name
    if (classic) {
        cmd <<= ' --classic'
    }
    if (channel != '') {
        cmd <<= " --channel=${channel}"
    }
    s(cmd)
}

void removeContainer() {
    // sudo -u $LXC_USER -H -E -- 
    sh '''
        lxc delete $BUILD_TAG --force
    '''
}

return this
