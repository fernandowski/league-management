module.exports = {
  preset: 'jest-expo',
  testMatch: ['**/__tests__/**/*.ts?(x)'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/$1',
  },
};
