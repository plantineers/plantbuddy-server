module.exports = {
    rules: {
        'body-leading-blank': [2, 'always'],
        'body-max-line-length': [2, 'always', 100],
        'header-max-length': [2, 'always', 100],
        'scope-case': [2, 'always', 'lower-case'],
        'subject-case': [
            2,
            'never',
            ['start-case', 'pascal-case', 'upper-case'],
        ],
        'subject-empty': [2, 'never'],
        'subject-full-stop': [2, 'never', '.'],
        'type-enum': [
            2,
            'always',
            [
                'chore',
                'build',
                'ci',
                'docs',
                'feat',
                'feat!',
                'fix',
                'perf',
                'refactor',
                'test'
            ],
        ],
        'type-case': [2, 'always', 'lower-case'],
        'type-empty': [2, 'never'],
        'signed-off-by': [2, 'always']
    }
};
