'use client';

import React, { useEffect } from 'react';
import { useGameStore } from '@/store/useGameStore';
import { useGameLoop } from '@/hooks/useGameLoop';

export default function GameDashboard() {
  // Start the Game Loop
  useGameLoop();

  const { character, inventory, syncWithServer, isSyncing } = useGameStore();

  // Initialize character if null (Mock login)
  useEffect(() => {
    if (!character) {
      useGameStore.getState().setCharacter({
        id: 'hero-123',
        name: 'Legendary Hero',
        hp: 100,
        maxHp: 100,
        gold: 0,
        level: 1,
      });
    }
  }, [character]);

  if (!character) return <div className="p-8 text-white">Loading Character...</div>;

  const hpPercentage = (character.hp / character.maxHp) * 100;

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 p-6 font-sans">
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Header Stats */}
        <div className="flex justify-between items-center bg-slate-800 p-4 rounded-xl border border-slate-700 shadow-xl">
          <div>
            <h1 className="text-2xl font-bold text-yellow-500">{character.name}</h1>
            <p className="text-slate-400 text-sm">Level {character.level} Adventurer</p>
          </div>
          <div className="text-right">
            <p className="text-xl font-mono font-bold text-yellow-400">{character.gold} GOLD</p>
            <button 
              onClick={() => syncWithServer()}
              disabled={isSyncing}
              className={`mt-2 px-4 py-1 rounded text-xs font-bold uppercase tracking-wider transition ${
                isSyncing ? 'bg-slate-600' : 'bg-blue-600 hover:bg-blue-500'
              }`}
            >
              {isSyncing ? 'Syncing...' : 'Manual Sync'}
            </button>
          </div>
        </div>

        {/* Health Bar */}
        <div className="space-y-2">
          <div className="flex justify-between text-sm font-bold">
            <span>HP</span>
            <span>{Math.floor(character.hp)} / {character.maxHp}</span>
          </div>
          <div className="w-full h-4 bg-slate-700 rounded-full overflow-hidden border border-slate-600">
            <div 
              className="h-full bg-red-500 transition-all duration-300 ease-out"
              style={{ width: `${hpPercentage}%` }}
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Battle Log Simulation Area */}
          <div className="bg-black/40 p-4 rounded-xl border border-slate-800 h-64 overflow-y-auto space-y-2 text-sm font-mono">
            <p className="text-green-400">[System] Welcome back, {character.name}!</p>
            <p className="text-slate-500">[System] Game loop started at 60fps...</p>
            <p className="text-blue-400">[Combat] Slime defeated! (Server-simulated)</p>
            <p className="text-slate-500">[System] Next auto-sync in 60s</p>
          </div>

          {/* Inventory Summary */}
          <div className="bg-slate-800 p-4 rounded-xl border border-slate-700">
            <h2 className="text-lg font-bold mb-4 border-b border-slate-700 pb-2">Inventory ({inventory.length}/100)</h2>
            <div className="grid grid-cols-4 gap-2">
              {inventory.slice(0, 8).map((item) => (
                <div key={item.id} className="aspect-square bg-slate-900 rounded border border-slate-700 flex items-center justify-center text-[10px] text-center p-1 group relative">
                  <span className="text-slate-300">{item.templateId.replace('drop_', '')}</span>
                  <div className="absolute inset-0 bg-blue-500/20 opacity-0 group-hover:opacity-100 transition-opacity rounded" />
                </div>
              ))}
              {inventory.length === 0 && <p className="col-span-4 text-slate-500 text-sm italic">Empty...</p>}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
