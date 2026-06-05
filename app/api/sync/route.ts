import { NextResponse } from 'next/server';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';

// Load gRPC definitions
const PROTO_PATH = path.join(process.cwd(), 'game.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});
const gameProto: any = grpc.loadPackageDefinition(packageDefinition).game;

// Create gRPC client
const client = new gameProto.GameService(
  'localhost:50051',
  grpc.credentials.createInsecure()
);

export async function POST(req: Request) {
  const body = await req.json();
  const { characterId, lastSync } = body;

  return new Promise((resolve) => {
    client.SyncGameState(
      {
        character_id: characterId,
        client_timestamp: Math.floor(lastSync / 1000),
      },
      (err: any, response: any) => {
        if (err) {
          console.error('gRPC Error:', err);
          resolve(NextResponse.json({ error: 'Failed to sync with game server' }, { status: 500 }));
        } else {
          // Map snake_case from gRPC to camelCase for Frontend
          const formattedResponse = {
            character: {
              id: response.state.character.id,
              name: response.state.character.name,
              hp: response.state.character.resources.current_hp,
              maxHp: response.state.character.resources.max_hp,
              gold: 150, // Mocked for now, usually from character stats/inventory
              level: response.state.character.level,
            },
            inventory: response.state.inventory.map((item: any) => ({
              id: item.id,
              templateId: item.template_id,
              level: item.level,
            })),
          };
          resolve(NextResponse.json(formattedResponse));
        }
      }
    );
  });
}
