/**
 * DO NOT CHANGE. This file is being managed from a central repository
 * To know more simply visit https://github.com/honestbank/.github/blob/main/docs/about.md
 */

class SemanticReleaseError extends Error {
    constructor(message, code, details) {
        super();
        Error.captureStackTrace(this, this.constructor);
        this.name = "SemanticReleaseError"
        this.details = details;
        this.code = code;
        this.semanticRelease = true;
    }
}

const getGolangVersion = () => `${process.env.GO_VERSION}` || '1.22'
const getTestDockerImageName = () => `asia.gcr.io/${process.env.DOCKER_REPO_TEST}`
const getProdDockerImageName = () => `asia.gcr.io/${process.env.DOCKER_REPO_PROD}`

module.exports = {
    branches: [{name: 'main'}, {name: 'master'}],
    plugins: [
        ["@semantic-release/commit-analyzer", {
            "preset": "angular",
            "releaseRules": [
                {type: 'feat', release: 'minor'},
                {type: 'fix', release: 'patch'},
                {type: 'perf', release: 'patch'},
                {type: 'docs', release: 'patch'},
                {type: 'refactor', release: 'patch'},
                {type: 'style', release: 'patch'},
                {type: 'chore', release: 'patch'}
            ]
        }],
        "@semantic-release/release-notes-generator",
        "@semantic-release/github"
    ],
    verifyConditions: [
        () => {
            if (!process.env.DOCKER_REPO_TEST) {
                throw new SemanticReleaseError(
                    "No DOCKER_REPO_TEST specified",
                    "ENODOCKER_REPO_TEST",
                    "Please make sure you're logged in to docker and a repo is available to push to"
                );
            }
            if (!process.env.DOCKER_REPO_PROD) {
                throw new SemanticReleaseError(
                    "No DOCKER_REPO_PROD specified",
                    "ENODOCKER_REPO_PROD",
                    "Please make sure you're logged in to docker and a repo is available to push to"
                );
            }
        },
        "@semantic-release/github"
    ],
    prepare: [
        {
            path: "@semantic-release/exec",
            cmd: `DOCKER_BUILDKIT=1 docker build . --build-arg GO_VERSION=${getGolangVersion()} --build-arg VERSION=\${nextRelease.version} -t ${getTestDockerImageName()}:\${nextRelease.version}`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker tag ${getTestDockerImageName()}:\${nextRelease.version} ${getTestDockerImageName()}:latest`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker tag ${getTestDockerImageName()}:\${nextRelease.version} ${getProdDockerImageName()}:\${nextRelease.version}`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker tag ${getTestDockerImageName()}:\${nextRelease.version} ${getProdDockerImageName()}:latest`
        }
    ],
    publish: [
        {
            path: "@semantic-release/exec",
            cmd: `docker --config /tmp/docker-test.json push ${getTestDockerImageName()}:\${nextRelease.version}`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker --config /tmp/docker-test.json push ${getTestDockerImageName()}:latest`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker --config /tmp/docker-prod.json push ${getProdDockerImageName()}:\${nextRelease.version}`
        },
        {
            path: "@semantic-release/exec",
            cmd: `docker --config /tmp/docker-prod.json push ${getProdDockerImageName()}:latest`
        },
        "@semantic-release/github"
    ],
};
