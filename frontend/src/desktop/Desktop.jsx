import React from 'react';
import Taskbar from './Taskbar';
import WorkspaceManager from './WorkspaceManager';
import '../styles/desktop.css';

export default function Desktop({ windows, layout, onLayoutChange }) {
  return (
    <div className="desktop-container">
      <WorkspaceManager windows={windows} layout={layout} onLayoutChange={onLayoutChange} />
      <Taskbar windows={windows} />
    </div>
  );
}