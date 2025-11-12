export const saveLayout = (layout) => {
  localStorage.setItem('aether-layout', JSON.stringify(layout));
};

export const restoreLayout = (callback) => {
  const savedLayout = localStorage.getItem('aether-layout');
  if (savedLayout) {
    callback(JSON.parse(savedLayout));
  }
};