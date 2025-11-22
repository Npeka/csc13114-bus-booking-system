import { NextResponse } from "next/server";
import { cookies } from "next/headers";

export async function POST() {
  try {
    // Get refresh token from httpOnly cookie before clearing
    const cookieStore = await cookies();
    const refreshToken = cookieStore.get("refresh_token")?.value;

    // Clear httpOnly cookie
    cookieStore.delete("refresh_token");

    return NextResponse.json({ success: true, refresh_token: refreshToken });
  } catch (error) {
    console.error("[Clear Token] Error clearing cookie:", error);
    return NextResponse.json(
      { error: "Failed to clear token" },
      { status: 500 },
    );
  }
}
