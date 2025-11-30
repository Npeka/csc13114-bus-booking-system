/**
 * Hook to get and check current user role
 */

import { useAuthStore } from "@/lib/stores/auth-store";
import { isAdmin, isSupport, isPassenger } from "./roles";

export function useRole() {
  const user = useAuthStore((state) => state.user);
  const userRole = user?.role ?? 0;

  return {
    role: userRole,
    isAdmin: isAdmin(userRole),
    isSupport: isSupport(userRole),
    isPassenger: isPassenger(userRole),
  };
}
