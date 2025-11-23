import { AuthenticatedTemplate } from "@/components/auth/authenticated-template";

/**
 * Shared template for all authenticated routes
 * Re-renders on navigation to check auth status
 */
export default function AuthTemplate({
  children,
}: {
  children: React.ReactNode;
}) {
  return <AuthenticatedTemplate>{children}</AuthenticatedTemplate>;
}
