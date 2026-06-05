import { useEffect, useRef } from 'react';
import { useGameStore } from '../store/useGameStore';

export const useGameLoop = () => {
  const updateHp = useGameStore((state) => state.updateHp);
  const addGold = useGameStore((state) => state.addGold);
  const syncWithServer = useGameStore((state) => state.syncWithServer);
  
  const requestRef = useRef<number>();
  const lastUpdateRef = useRef<number>(performance.now());
  const lastSyncRef = useRef<number>(performance.now());

  const animate = (time: number) => {
    const deltaTime = time - lastUpdateRef.current;

    // 1. Visual Simulation (Run every frame for 60fps)
    // Here we simulate small HP regen or monster damage for visual feedback
    // Note: Real state changes still get corrected by the sync.
    if (deltaTime >= 16.67) { // ~60fps
      // Example: Small HP regen visual simulation
      updateHp(0.01 * (deltaTime / 16.67));
      lastUpdateRef.current = time;
    }

    // 2. Scheduled Sync (Every 60 seconds)
    if (time - lastSyncRef.current >= 60000) {
      syncWithServer();
      lastSyncRef.current = time;
    }

    requestRef.current = requestAnimationFrame(animate);
  };

  useEffect(() => {
    requestRef.current = requestAnimationFrame(animate);
    return () => {
      if (requestRef.current) {
        cancelAnimationFrame(requestRef.current);
      }
    };
  }, []);
};
