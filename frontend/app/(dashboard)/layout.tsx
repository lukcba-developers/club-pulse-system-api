import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { AuthProvider } from "@/context/auth-context";
import { Toaster } from "@/components/ui/use-toast";

import { BrandProvider } from "@/components/providers/BrandProvider";

export default function Layout({ children }: { children: React.ReactNode }) {
    return (
        <AuthProvider>
            <BrandProvider>
                <DashboardLayout>{children}</DashboardLayout>
            </BrandProvider>
            <Toaster />
        </AuthProvider>
    );
}

