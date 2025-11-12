import React, { useState, useEffect } from 'react';
import Taskbar from './Taskbar';
import WorkspaceManager from './WorkspaceManager';
import Dock from './Dock';
import DesktopShortcut from './DesktopShortcut';
import '../styles/desktop.css';

export default function Desktop({ windows, layout, onLayoutChange, onOpenWindow, onCloseWindow, onMinimizeWindow, onFocusWindow }) {
  const [shortcuts, setShortcuts] = useState([]);

  useEffect(() => {
    fetch('/v1/apps')
      .then(response => response.json())
      .then(apps => {
        const desktopShortcuts = apps.map(app => ({
          id: app.id,
          name: app.name,
          icon: app.icon,
          window: app.window,
        }));
        setShortcuts(desktopShortcuts);
      });
  }, []);

  return (
    <div className="desktop-container">
      <WorkspaceManager windows={windows} layout={layout} onLayoutChange={onLayoutChange} onCloseWindow={onCloseWindow} onMinimizeWindow={onMinimizeWindow} onFocusWindow={onFocusWindow} />
      {shortcuts.map(shortcut => (
        <DesktopShortcut key={shortcut.id} shortcut={shortcut} onOpenWindow={onOpenWindow} />
      ))}
      <Taskbar windows={windows} onOpenWindow={onOpenWindow} />
      <Dock onOpenWindow={onOpenWindow} apps={shortcuts} />
    </div>
  );
}
