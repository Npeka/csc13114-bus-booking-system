/**
 * User Role System
 * Backend role model: Bit-flag based (Passenger=1, Admin=2, Support=8)
 * Frontend role model: Priority-based (highest role takes precedence)
 */

export enum Role {
  GUEST = 0,
  PASSENGER = 2,
  ADMIN = 4,
}

export const ROLE_NAMES: Record<Role, string> = {
  [Role.GUEST]: "Guest",
  [Role.PASSENGER]: "Passenger",
  [Role.ADMIN]: "Admin",
};

export const ROLE_LABELS: Record<Role, string> = {
  [Role.GUEST]: "Khách vãng lai",
  [Role.PASSENGER]: "Hành khách",
  [Role.ADMIN]: "Quản trị viên",
};

/**
 * Check if a user has a specific role (bit-flag check)
 */
export function hasRole(userRole: number, role: Role): boolean {
  if (!userRole) return false;
  return (userRole & role) === role;
}

/**
 * Check if user is admin
 */
export function isAdmin(userRole: number): boolean {
  return hasRole(userRole, Role.ADMIN);
}

/**
 * Check if user is passenger (basic customer)
 */
export function isPassenger(userRole: number): boolean {
  return hasRole(userRole, Role.PASSENGER);
}

/**
 * Get list of roles for a user (can have multiple due to bit-flag)
 */
export function getUserRoles(userRole: number): Role[] {
  const roles: Role[] = [];
  if (userRole === Role.GUEST) {
    roles.push(Role.GUEST);
  } else {
    if (isPassenger(userRole)) roles.push(Role.PASSENGER);
    if (isAdmin(userRole)) roles.push(Role.ADMIN);
  }
  return roles;
}

/**
 * Get role labels for a user
 */
export function getUserRoleLabels(userRole: number): string[] {
  return getUserRoles(userRole).map((role) => ROLE_LABELS[role]);
}

/**
 * Get primary role name (for display)
 */
export function getRoleName(userRole: number): string {
  // Return the highest priority role name
  if (isAdmin(userRole)) return ROLE_NAMES[Role.ADMIN];
  return ROLE_NAMES[Role.PASSENGER];
}

/**
 * Get primary role label in Vietnamese (for display)
 */
export function getRoleLabel(userRole: number): string {
  // Return the highest priority role label (Vietnamese)
  if (isAdmin(userRole)) return ROLE_LABELS[Role.ADMIN];
  return ROLE_LABELS[Role.PASSENGER];
}
