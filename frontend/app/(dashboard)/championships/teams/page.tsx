"use client"

import { TeamManagement } from "@/components/championship/TeamManagement"

export default function TeamsPage() {


    return (
        <div className="container mx-auto py-6">
            <h1 className="text-3xl font-bold mb-6">Equipos</h1>
            <div className="max-w-4xl">
                <TeamManagement />
            </div>
        </div>
    )
}
