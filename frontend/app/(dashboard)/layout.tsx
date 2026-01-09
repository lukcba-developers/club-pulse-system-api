import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { AuthProvider } from "@/context/auth-context";
import { Toaster } from "@/components/ui/use-toast";

export default function Layout({ children }: { children: React.ReactNode }) {
    return (
        <AuthProvider>
            <DashboardLayout>{children}</DashboardLayout>
            <Toaster />
        </AuthProvider>
    );
}

