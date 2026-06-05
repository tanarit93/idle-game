import { NextResponse } from 'next/server';

export async function POST(req: Request) {
  try {
    const { characterId, lastSync } = await req.json();

    // In a production environment, this would call the gRPC server at localhost:50051
    // For this local setup, since the Go server exits after one run for the demo,
    // we would normally keep it running as a service.
    
    // Let's assume the Go server is running in 'service' mode in Docker.
    // For now, we provide the mock structure that matches the Go SyncResponse.
    
    return NextResponse.json({
      character: {
        id: characterId,
        name: "Legendary Hero",
        level: 5,
        hp: 100,
        maxHp: 100,
        gold: 540,
        experience: 1200
      },
      inventory: [
        { id: "1", templateId: "iron_scrap", level: 1 },
        { id: "2", templateId: "fire_skill_gem", level: 1 }
      ]
    });
  } catch (error) {
    return NextResponse.json({ error: 'Failed to communicate with Game Engine' }, { status: 500 });
  }
}
