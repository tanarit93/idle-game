import { create } from 'zustand';

export interface Item {
  id: string;
  templateId: string;
  type: string;
  tier: number;
  slottedGems?: Item[];
}

export interface Character {
  id: string;
  name: string;
  hp: number;
  maxHp: number;
  gold: number;
  level: number;
}

interface GameState {
  character: Character | null;
  inventory: Item[];
  lastSync: number;
  isSyncing: boolean;

  // Actions
  setCharacter: (char: Character) => void;
  updateHp: (amount: number) => void;
  addGold: (amount: number) => void;
  syncWithServer: () => Promise<void>;
}

export const useGameStore = create<GameState>((set, get) => ({
  character: null,
  inventory: [],
  lastSync: Date.now(),
  isSyncing: false,

  setCharacter: (char) => set({ character: char }),

  updateHp: (amount) => set((state) => {
    if (!state.character) return state;
    const newHp = Math.max(0, Math.min(state.character.maxHp, state.character.hp + amount));
    return { character: { ...state.character, hp: newHp } };
  }),

  addGold: (amount) => set((state) => {
    if (!state.character) return state;
    return { character: { ...state.character, gold: state.character.gold + amount } };
  }),

  syncWithServer: async () => {
    const { character, lastSync, isSyncing } = get();
    if (!character || isSyncing) return;

    set({ isSyncing: true });

    try {
      const response = await fetch('/api/sync', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          characterId: character.id,
          lastSync,
        }),
      });

      if (response.ok) {
        const authoritativeState = await response.json();
        // Forcefully overwrite local state with server truth
        set({
          character: authoritativeState.character,
          inventory: authoritativeState.inventory,
          lastSync: Date.now(),
        });
      }
    } catch (error) {
      console.error('Failed to sync game state:', error);
    } finally {
      set({ isSyncing: false });
    }
  },
}));
