"use client"

import { useState, useEffect } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { SessionList } from "@/components/session-list"
import { FamilyList } from "@/components/profile/family-list"
import { BillingSection } from "@/components/profile/billing-section"
import { GamificationStats } from "@/components/profile/gamification-stats"
import { HealthSection } from "@/components/profile/health-section"
import { DocumentUpload, UserDocument } from "@/components/DocumentUpload"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { CheckCircle, XCircle } from "lucide-react"

export default function ProfilePage() {
    const [documents, setDocuments] = useState<UserDocument[]>([])
    const [eligibility, setEligibility] = useState<{
        is_eligible: boolean
        issues: string[]
        has_dni: boolean
        has_emmac: boolean
    } | null>(null)
    const [userId, setUserId] = useState<string>("")

    useEffect(() => {
        // TODO: Obtener userId del contexto de autenticación
        const mockUserId = "current-user-id"
        setUserId(mockUserId)
        fetchDocuments(mockUserId)
        fetchEligibility(mockUserId)
    }, [])

    const fetchDocuments = async (uid: string) => {
        try {
            const response = await fetch(`/api/users/${uid}/documents`, {
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            })
            if (response.ok) {
                const data = await response.json()
                setDocuments(data)
            }
        } catch (error) {
            console.error("Error fetching documents:", error)
        }
    }

    const fetchEligibility = async (uid: string) => {
        try {
            const response = await fetch(`/api/users/${uid}/eligibility`, {
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            })
            if (response.ok) {
                const data = await response.json()
                setEligibility(data)
            }
        } catch (error) {
            console.error("Error fetching eligibility:", error)
        }
    }

    const handleUploadSuccess = () => {
        fetchDocuments(userId)
        fetchEligibility(userId)
    }

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-bold tracking-tight">Perfil & Estadísticas</h1>
            </div>

            {/* Eligibility Status Card */}
            {eligibility && (
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            Estado de Elegibilidad
                            {eligibility.is_eligible ? (
                                <Badge variant="success" className="ml-2">
                                    <CheckCircle className="h-4 w-4 mr-1" />
                                    Habilitado
                                </Badge>
                            ) : (
                                <Badge variant="destructive" className="ml-2">
                                    <XCircle className="h-4 w-4 mr-1" />
                                    Inhabilitado
                                </Badge>
                            )}
                        </CardTitle>
                        <CardDescription>
                            {eligibility.is_eligible
                                ? "Tienes toda la documentación necesaria para participar"
                                : "Necesitas completar tu documentación para poder participar"}
                        </CardDescription>
                    </CardHeader>
                    {!eligibility.is_eligible && eligibility.issues.length > 0 && (
                        <CardContent>
                            <div className="space-y-2">
                                <p className="text-sm font-medium">Problemas detectados:</p>
                                <ul className="list-disc list-inside space-y-1">
                                    {eligibility.issues.map((issue, index) => (
                                        <li key={index} className="text-sm text-muted-foreground">
                                            {issue}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        </CardContent>
                    )}
                </Card>
            )}

            <Tabs defaultValue="overview" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="overview">General</TabsTrigger>
                    <TabsTrigger value="documents">Documentación</TabsTrigger>
                    <TabsTrigger value="billing">Facturación</TabsTrigger>
                    <TabsTrigger value="family">Familia</TabsTrigger>
                    <TabsTrigger value="sessions">Sesiones</TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="space-y-6">
                    <GamificationStats />
                    <HealthSection />
                </TabsContent>

                <TabsContent value="documents" className="space-y-6">
                    <DocumentUpload
                        userId={userId}
                        documents={documents}
                        onUploadSuccess={handleUploadSuccess}
                    />
                </TabsContent>

                <TabsContent value="billing" className="space-y-6">
                    <BillingSection />
                </TabsContent>

                <TabsContent value="family" className="space-y-6">
                    <FamilyList />
                </TabsContent>

                <TabsContent value="sessions" className="space-y-6">
                    <SessionList />
                </TabsContent>
            </Tabs>
        </div>
    )
}
