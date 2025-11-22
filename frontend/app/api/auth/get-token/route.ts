import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function POST() {
  try {
    // Get refresh token from httpOnly cookie
    const cookieStore = await cookies();
    const refreshToken = cookieStore.get("refresh_token")?.value;

    return NextResponse.json({ refresh_token: refreshToken || null });
  } catch (error) {
    console.error("[Get Token] Error reading token:", error);
    return NextResponse.json(
      { error: "Failed to get token", refresh_token: null },
      { status: 500 },
    );
  }
}
