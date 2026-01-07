import { Skeleton } from "@/components/ui/skeleton"

export default function DashboardLoading() {
    return (
        <div className="max-w-7xl mx-auto p-4 sm:p-6 lg:p-8 animate-in fade-in duration-500">
            {/* Header Skeleton */}
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-8 gap-4">
                <div className="space-y-2">
                    <Skeleton className="h-8 w-48" />
                    <Skeleton className="h-4 w-32" />
                </div>
                <Skeleton className="h-10 w-40" />
            </div>

            {/* Content Grid Skeleton */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">

                {/* Main Column */}
                <div className="lg:col-span-2 space-y-8">
                    {/* Alert Skeleton */}
                    <Skeleton className="h-24 w-full rounded-xl" />

                    {/* Search Skeleton */}
                    <div className="flex justify-center">
                        <Skeleton className="h-12 w-full max-w-2xl rounded-lg" />
                    </div>

                    {/* Facilities Grid Skeleton */}
                    <div className="space-y-4">
                        <Skeleton className="h-6 w-40 mb-4" />
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <Skeleton className="h-64 rounded-xl" />
                            <Skeleton className="h-64 rounded-xl" />
                            <Skeleton className="h-64 rounded-xl" />
                            <Skeleton className="h-64 rounded-xl" />
                        </div>
                    </div>
                </div>

                {/* Sidebar Skeleton */}
                <div className="lg:col-span-1">
                    <Skeleton className="h-96 rounded-xl" />
                </div>

            </div>
        </div>
    )
}
