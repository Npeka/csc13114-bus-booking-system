export default function PageHeaderLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="mb-8 flex items-center justify-between">{children}</div>
  );
}
