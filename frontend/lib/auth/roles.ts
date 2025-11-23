/**
 * Role-based access control utilities
 * Backend role model: Bit-flag based (Passenger=1, Admin=2, Operator=4, Support=8)
 */

export enum Role {
  PASSENGER = 1,
  ADMIN = 2,
  OPERATOR = 4,
  SUPPORT = 8,
}

export const ROLE_NAMES: Record<Role, string> = {
  [Role.PASSENGER]: "Passenger",
  [Role.ADMIN]: "Admin",
  [Role.OPERATOR]: "Operator",
  [Role.SUPPORT]: "Support",
};

export const ROLE_LABELS: Record<Role, string> = {
  [Role.PASSENGER]: "Khách hàng",
  [Role.ADMIN]: "Quản trị viên",
  [Role.OPERATOR]: "Nhà điều hành",
  [Role.SUPPORT]: "Hỗ trợ",
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
 * Check if user is operator
 */
export function isOperator(userRole: number): boolean {
  return hasRole(userRole, Role.OPERATOR);
}

/**
 * Check if user is support
 */
export function isSupport(userRole: number): boolean {
  return hasRole(userRole, Role.SUPPORT);
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
  if (isPassenger(userRole)) roles.push(Role.PASSENGER);
  if (isAdmin(userRole)) roles.push(Role.ADMIN);
  if (isOperator(userRole)) roles.push(Role.OPERATOR);
  if (isSupport(userRole)) roles.push(Role.SUPPORT);
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
export function getPrimaryRoleName(userRole: number): string {
  if (isAdmin(userRole)) return ROLE_NAMES[Role.ADMIN];
  if (isOperator(userRole)) return ROLE_NAMES[Role.OPERATOR];
  if (isSupport(userRole)) return ROLE_NAMES[Role.SUPPORT];
  return ROLE_NAMES[Role.PASSENGER];
}

/**
 * Get primary role label in Vietnamese (for display)
 */
export function getPrimaryRoleLabel(userRole: number): string {
  if (isAdmin(userRole)) return ROLE_LABELS[Role.ADMIN];
  if (isOperator(userRole)) return ROLE_LABELS[Role.OPERATOR];
  if (isSupport(userRole)) return ROLE_LABELS[Role.SUPPORT];
  return ROLE_LABELS[Role.PASSENGER];
}
