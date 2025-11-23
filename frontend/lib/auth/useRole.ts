/**
 * Hook to get and check current user role
 */

import { useAuthStore } from "@/lib/stores/auth-store";
import { isAdmin, isOperator, isSupport, isPassenger } from "./roles";

export function useRole() {
  const user = useAuthStore((state) => state.user);
  const userRole = user?.role ?? 0;

  return {
    role: userRole,
    isAdmin: isAdmin(userRole),
    isOperator: isOperator(userRole),
    isSupport: isSupport(userRole),
    isPassenger: isPassenger(userRole),
  };
}
