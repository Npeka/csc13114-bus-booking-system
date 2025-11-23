import { render, screen } from "@/lib/test-utils";
import { RoleBadge } from "@/components/auth/role-badge";
import { Role } from "@/lib/auth/roles";

describe("RoleBadge component", () => {
  it("should render role badge with admin label", () => {
    render(<RoleBadge userRole={Role.ADMIN} />);
    expect(screen.getByText("Quản trị viên")).toBeInTheDocument();
  });

  it("should render role badge with operator label", () => {
    render(<RoleBadge userRole={Role.OPERATOR} />);
    expect(screen.getByText("Nhà điều hành")).toBeInTheDocument();
  });

  it("should render role badge with passenger label", () => {
    render(<RoleBadge userRole={Role.PASSENGER} />);
    expect(screen.getByText("Khách hàng")).toBeInTheDocument();
  });

  it("should render nothing for zero role", () => {
    const { container } = render(<RoleBadge userRole={0} />);
    // The container has the theme script as firstChild, so check that the Badge doesn't exist
    expect(
      screen.queryByText(/Quản trị viên|Nhà điều hành|Khách hàng|Hỗ trợ/),
    ).not.toBeInTheDocument();
  });

  it("should prioritize admin in multi-role", () => {
    const multiRole = Role.ADMIN | Role.PASSENGER;
    render(<RoleBadge userRole={multiRole} />);
    expect(screen.getByText("Quản trị viên")).toBeInTheDocument();
  });

  it("should accept custom className", () => {
    render(<RoleBadge userRole={Role.ADMIN} className="custom-class" />);
    const badge = screen.getByText("Quản trị viên").closest("span");
    expect(badge).toHaveClass("custom-class");
  });
});
