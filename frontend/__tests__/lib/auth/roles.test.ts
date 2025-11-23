import {
  hasRole,
  isAdmin,
  isOperator,
  isSupport,
  isPassenger,
  Role,
  getPrimaryRoleName,
  getPrimaryRoleLabel,
  getUserRoles,
} from "@/lib/auth/roles";

describe("Role utilities", () => {
  describe("hasRole", () => {
    it("should correctly identify admin role", () => {
      const adminRole = Role.ADMIN; // 2
      expect(hasRole(adminRole, Role.ADMIN)).toBe(true);
      expect(hasRole(adminRole, Role.PASSENGER)).toBe(false);
    });

    it("should handle bit-flag roles", () => {
      const multiRole = Role.ADMIN | Role.PASSENGER; // 3
      expect(hasRole(multiRole, Role.ADMIN)).toBe(true);
      expect(hasRole(multiRole, Role.PASSENGER)).toBe(true);
      expect(hasRole(multiRole, Role.OPERATOR)).toBe(false);
    });

    it("should return false for zero role", () => {
      expect(hasRole(0, Role.ADMIN)).toBe(false);
    });
  });

  describe("isAdmin", () => {
    it("should return true for admin", () => {
      expect(isAdmin(Role.ADMIN)).toBe(true);
    });

    it("should return false for passenger", () => {
      expect(isAdmin(Role.PASSENGER)).toBe(false);
    });

    it("should return true for admin in multi-role", () => {
      const multiRole = Role.ADMIN | Role.PASSENGER;
      expect(isAdmin(multiRole)).toBe(true);
    });
  });

  describe("isOperator", () => {
    it("should return true for operator", () => {
      expect(isOperator(Role.OPERATOR)).toBe(true);
    });

    it("should return false for passenger", () => {
      expect(isOperator(Role.PASSENGER)).toBe(false);
    });
  });

  describe("isSupport", () => {
    it("should return true for support", () => {
      expect(isSupport(Role.SUPPORT)).toBe(true);
    });

    it("should return false for passenger", () => {
      expect(isSupport(Role.PASSENGER)).toBe(false);
    });
  });

  describe("isPassenger", () => {
    it("should return true for passenger", () => {
      expect(isPassenger(Role.PASSENGER)).toBe(true);
    });

    it("should return false for admin", () => {
      expect(isPassenger(Role.ADMIN)).toBe(false);
    });
  });

  describe("getPrimaryRoleName", () => {
    it("should return admin name for admin role", () => {
      expect(getPrimaryRoleName(Role.ADMIN)).toBe("Admin");
    });

    it("should return passenger name for passenger role", () => {
      expect(getPrimaryRoleName(Role.PASSENGER)).toBe("Passenger");
    });

    it("should prioritize admin in multi-role", () => {
      const multiRole = Role.ADMIN | Role.PASSENGER;
      expect(getPrimaryRoleName(multiRole)).toBe("Admin");
    });
  });

  describe("getPrimaryRoleLabel", () => {
    it("should return Vietnamese admin label", () => {
      expect(getPrimaryRoleLabel(Role.ADMIN)).toBe("Quản trị viên");
    });

    it("should return Vietnamese passenger label", () => {
      expect(getPrimaryRoleLabel(Role.PASSENGER)).toBe("Khách hàng");
    });
  });

  describe("getUserRoles", () => {
    it("should return array of roles for user", () => {
      const multiRole = Role.ADMIN | Role.PASSENGER;
      const roles = getUserRoles(multiRole);
      expect(roles).toContain(Role.ADMIN);
      expect(roles).toContain(Role.PASSENGER);
      expect(roles.length).toBe(2);
    });

    it("should return single role array", () => {
      const roles = getUserRoles(Role.ADMIN);
      expect(roles).toEqual([Role.ADMIN]);
    });

    it("should return empty array for zero role", () => {
      const roles = getUserRoles(0);
      expect(roles).toEqual([]);
    });
  });
});
