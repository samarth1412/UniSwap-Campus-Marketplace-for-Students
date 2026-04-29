import '@testing-library/jest-dom';

if (!URL.createObjectURL) {
  Object.defineProperty(URL, 'createObjectURL', {
    value: () => 'blob:mock-preview',
    writable: true,
  });
}

if (!URL.revokeObjectURL) {
  Object.defineProperty(URL, 'revokeObjectURL', {
    value: () => undefined,
    writable: true,
  });
}
