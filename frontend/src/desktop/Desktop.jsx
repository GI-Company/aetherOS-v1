import React from 'react';
import WorkspaceManager from './WorkspaceManager';
import DesktopShortcut from './DesktopShortcut';
import '../styles/desktop.css';

export default function Desktop({ windows, layout, onLayoutChange, onOpenWindow, onCloseWindow, onMinimizeWindow, onFocusWindow, shortcuts }) {
  return (
    <div className="desktop-container">
      <WorkspaceManager windows={windows} layout={layout} onLayoutChange={onLayoutChange} onCloseWindow={onCloseWindow} onMinimizeWindow={onMinimizeWindow} onFocusWindow={onFocusWindow} />
      {shortcuts.map(shortcut => (
        <DesktopShortcut key={shortcut.id} shortcut={shortcut} onOpenWindow={onOpenWindow} />
      ))}
    </div>
  );
}
