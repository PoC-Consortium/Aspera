module.exports = {
    preset: "jest-preset-angular",
    roots: ['src'],
    setupTestFrameworkScriptFile: "./src/setup-jest.ts",
    moduleNameMapper: {
        '@app/(.*)': 'src/app/$1',
    }
}