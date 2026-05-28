import { useAuth } from "@/contexts/AuthContext";

interface CanProps {
  permission?: string;
  permissions?: string[];
  match?: "all" | "any";
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function Can({
  permission,
  permissions,
  match = "all",
  children,
  fallback = null,
}: CanProps) {
  const { can } = useAuth();

  const allPerms = [
    ...(permission ? [permission] : []),
    ...(permissions ?? []),
  ];

  if (allPerms.length === 0) return <>{children}</>;

  const granted =
    match === "any"
      ? allPerms.some((p) => can(p))
      : allPerms.every((p) => can(p));

  return granted ? <>{children}</> : <>{fallback}</>;
}