"use client"

import { useState } from "react"
import { Upload, FileText, CheckCircle, XCircle, Clock, AlertCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { useToast } from "@/hooks/use-toast"

export type DocumentType = "DNI_FRONT" | "DNI_BACK" | "EMMAC_MEDICAL" | "LEAGUE_FORM" | "INSURANCE"
export type DocumentStatus = "PENDING" | "VALID" | "REJECTED" | "EXPIRED"

export interface UserDocument {
    id: string
    type: DocumentType
    file_url: string
    status: DocumentStatus
    expiration_date?: string
    uploaded_at: string
    validated_at?: string
    rejection_notes?: string
}

interface DocumentUploadProps {
    userId: string
    documents: UserDocument[]
    onUploadSuccess?: () => void
}

const DOCUMENT_TYPES: Record<DocumentType, string> = {
    DNI_FRONT: "DNI (Frente)",
    DNI_BACK: "DNI (Dorso)",
    EMMAC_MEDICAL: "Apto Médico (EMMAC)",
    LEAGUE_FORM: "Formulario de Liga",
    INSURANCE: "Seguro"
}

const STATUS_CONFIG: Record<DocumentStatus, { label: string; icon: React.ReactNode; variant: "default" | "success" | "destructive" | "warning" }> = {
    PENDING: { label: "En revisión", icon: <Clock className="h-4 w-4" />, variant: "warning" },
    VALID: { label: "Válido", icon: <CheckCircle className="h-4 w-4" />, variant: "success" },
    REJECTED: { label: "Rechazado", icon: <XCircle className="h-4 w-4" />, variant: "destructive" },
    EXPIRED: { label: "Vencido", icon: <AlertCircle className="h-4 w-4" />, variant: "destructive" }
}

export function DocumentUpload({ userId, documents, onUploadSuccess }: DocumentUploadProps) {
    const [selectedType, setSelectedType] = useState<DocumentType>("DNI_FRONT")
    const [selectedFile, setSelectedFile] = useState<File | null>(null)
    const [expirationDate, setExpirationDate] = useState("")
    const [isUploading, setIsUploading] = useState(false)
    const { toast } = useToast()

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setSelectedFile(e.target.files[0])
        }
    }

    const handleUpload = async () => {
        if (!selectedFile) {
            toast({
                title: "Error",
                description: "Por favor selecciona un archivo",
                variant: "destructive"
            })
            return
        }

        setIsUploading(true)

        try {
            const formData = new FormData()
            formData.append("file", selectedFile)
            formData.append("type", selectedType)
            if (expirationDate) {
                formData.append("expiration_date", expirationDate)
            }

            const response = await fetch(`/api/users/${userId}/documents`, {
                method: "POST",
                body: formData,
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            })

            if (!response.ok) {
                throw new Error("Error al subir documento")
            }

            toast({
                title: "✅ Documento subido",
                description: "El documento ha sido enviado para revisión"
            })

            // Reset form
            setSelectedFile(null)
            setExpirationDate("")
            if (onUploadSuccess) {
                onUploadSuccess()
            }
        } catch {
            toast({
                title: "Error",
                description: "No se pudo subir el documento",
                variant: "destructive"
            })
        } finally {
            setIsUploading(false)
        }
    }

    const getDocumentByType = (type: DocumentType) => {
        return documents.find(doc => doc.type === type)
    }

    return (
        <div className="space-y-6">
            {/* Upload Form */}
            <Card>
                <CardHeader>
                    <CardTitle>Subir Documento</CardTitle>
                    <CardDescription>
                        Sube tus documentos para que sean validados por el administrador
                    </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="document-type">Tipo de Documento</Label>
                        <Select value={selectedType} onValueChange={(value) => setSelectedType(value as DocumentType)}>
                            <SelectTrigger id="document-type">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                {Object.entries(DOCUMENT_TYPES).map(([key, label]) => (
                                    <SelectItem key={key} value={key}>
                                        {label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="file">Archivo</Label>
                        <Input
                            id="file"
                            type="file"
                            accept="image/*,application/pdf"
                            onChange={handleFileChange}
                        />
                        {selectedFile && (
                            <p className="text-sm text-muted-foreground">
                                Archivo seleccionado: {selectedFile.name}
                            </p>
                        )}
                    </div>

                    {(selectedType === "EMMAC_MEDICAL" || selectedType === "INSURANCE") && (
                        <div className="space-y-2">
                            <Label htmlFor="expiration">Fecha de Vencimiento</Label>
                            <Input
                                id="expiration"
                                type="date"
                                value={expirationDate}
                                onChange={(e) => setExpirationDate(e.target.value)}
                            />
                        </div>
                    )}

                    <Button onClick={handleUpload} disabled={isUploading || !selectedFile} className="w-full">
                        <Upload className="mr-2 h-4 w-4" />
                        {isUploading ? "Subiendo..." : "Subir Documento"}
                    </Button>
                </CardContent>
            </Card>

            {/* Documents List */}
            <Card>
                <CardHeader>
                    <CardTitle>Mis Documentos</CardTitle>
                    <CardDescription>
                        Estado de tus documentos subidos
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <div className="space-y-3">
                        {Object.entries(DOCUMENT_TYPES).map(([type, label]) => {
                            const doc = getDocumentByType(type as DocumentType)
                            const status = doc?.status || null

                            return (
                                <div key={type} className="flex items-center justify-between p-3 border rounded-lg">
                                    <div className="flex items-center gap-3">
                                        <FileText className="h-5 w-5 text-muted-foreground" />
                                        <div>
                                            <p className="font-medium">{label}</p>
                                            {doc?.expiration_date && (
                                                <p className="text-sm text-muted-foreground">
                                                    Vence: {new Date(doc.expiration_date).toLocaleDateString()}
                                                </p>
                                            )}
                                            {doc?.rejection_notes && (
                                                <p className="text-sm text-destructive">
                                                    Motivo: {doc.rejection_notes}
                                                </p>
                                            )}
                                        </div>
                                    </div>

                                    {status ? (
                                        <Badge variant={STATUS_CONFIG[status].variant} className="flex items-center gap-1">
                                            {STATUS_CONFIG[status].icon}
                                            {STATUS_CONFIG[status].label}
                                        </Badge>
                                    ) : (
                                        <Badge variant="outline">Sin subir</Badge>
                                    )}
                                </div>
                            )
                        })}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
