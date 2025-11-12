import React from 'react';
import '../styles/desktopShortcut.css';

export default function DesktopShortcut({ shortcut, onOpenWindow }) {
  const handleDoubleClick = () => {
    onOpenWindow(shortcut.window);
  };

  return (
    <div
      className="desktop-shortcut"
      style={{ left: shortcut.x, top: shortcut.y }}
      onDoubleClick={handleDoubleClick}
    >
      <img src={shortcut.icon} alt={shortcut.name} />
      <span>{shortcut.name}</span>
    </div>
  );
}
