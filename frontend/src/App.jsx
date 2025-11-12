
import React, { useEffect, useState } from 'react';
import Desktop from './desktop/Desktop';
import { AIAgentWindow, ComputeWindow, NotificationWindow, MarketplaceWindow } from './windows';
import { restoreLayout, saveLayout } from './sdk/layoutManager';

export default function App() {
  const [layout, setLayout] = useState({});
  const [windows, setWindows] = useState([]);

  useEffect(() => {
    restoreLayout((savedLayout) => setLayout(savedLayout));
  }, []);

  const handleLayoutChange = (newLayout) => {
    setLayout(newLayout);
    saveLayout(newLayout);
  };

  const handleOpenWindow = (windowId) => {
    if (windows.find((w) => w.id === windowId)) {
      return;
    }

    let newWindow;
    switch (windowId) {
      case 'ai':
        newWindow = { id: 'ai', title: 'AI Agent', component: <AIAgentWindow /> };
        break;
      case 'compute':
        newWindow = { id: 'compute', title: 'Compute', component: <ComputeWindow /> };
        break;
      case 'notifications':
        newWindow = { id: 'notifications', title: 'Notifications', component: <NotificationWindow /> };
        break;
      case 'marketplace':
        newWindow = { id: 'marketplace', title: 'Marketplace', component: <MarketplaceWindow /> };
        break;
      default:
        return;
    }

    setWindows((prevWindows) => [...prevWindows, newWindow]);
  };

  const handleCloseWindow = (windowId) => {
    setWindows((prevWindows) => prevWindows.filter((w) => w.id !== windowId));
  };

  const handleMinimizeWindow = (windowId) => {
    setWindows((prevWindows) =>
      prevWindows.map((w) => (w.id === windowId ? { ...w, minimized: true } : w))
    );
  };
  
  const advanceWindows = () => {
    const newLayout = { ...layout };
    windows.forEach((w, i) => {
      newLayout[w.id] = {
        ...newLayout[w.id],
        x: 50 + i * 20,
        y: 50 + i * 20,
        zIndex: 1000 + i
      };
    });
    setLayout(newLayout);
  };

  const handleFocusWindow = (windowId) => {
    const newLayout = { ...layout };
    const maxZIndex = Math.max(...Object.values(newLayout).map(l => (l && l.zIndex) || 0));
    newLayout[windowId] = { ...newLayout[windowId], zIndex: maxZIndex + 1 };
    setLayout(newLayout);
  };


  return <Desktop windows={windows} layout={layout} onLayoutChange={handleLayoutChange} onOpenWindow={handleOpenWindow} onCloseWindow={handleCloseWindow} onMinimizeWindow={handleMinimizeWindow} advanceWindows={advanceWindows} onFocusWindow={handleFocusWindow}/>;
}
