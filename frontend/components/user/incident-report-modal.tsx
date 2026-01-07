"use client";

import { useState } from "react";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { userService } from "@/services/user-service";
import { useToast } from "@/components/ui/use-toast";
import { AlertTriangle } from "lucide-react";

export function IncidentReportModal() {
    const [open, setOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const { toast } = useToast();

    const [form, setForm] = useState({
        description: "",
        witnesses: "",
        action_taken: "",
        injured_user_id: "" // Optional
    });

    const handleSubmit = async () => {
        try {
            setLoading(true);
            await userService.logIncident({
                description: form.description,
                witnesses: form.witnesses,
                action_taken: form.action_taken,
                injured_user_id: form.injured_user_id || undefined
            });
            toast({
                title: "Incidente Registrado",
                description: "El reporte ha sido guardado en el sistema de seguridad.",
            });
            setOpen(false);
            setForm({ description: "", witnesses: "", action_taken: "", injured_user_id: "" });
        } catch (error) {
            console.error(error);
            toast({
                title: "Error",
                description: "No se pudo registrar el incidente.",
                variant: "destructive"
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button variant="destructive" className="w-full justify-start bg-red-600 hover:bg-red-700 text-white">
                    <AlertTriangle className="mr-2 h-4 w-4" />
                    Reportar Incidente
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2 text-red-600">
                        <AlertTriangle className="h-5 w-5" />
                        Reporte de Incidente / Accidente
                    </DialogTitle>
                    <DialogDescription>
                        Utilice este formulario para registrar accidentes, lesiones o problemas de seguridad en el club.
                    </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                    <div className="grid gap-2">
                        <Label htmlFor="description">Descripción del Hecho *</Label>
                        <Textarea
                            id="description"
                            placeholder="Detalle qué ocurrió, cuándo y dónde..."
                            value={form.description}
                            onChange={(e) => setForm({ ...form, description: e.target.value })}
                        />
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="action">Acción Tomada</Label>
                        <Input
                            id="action"
                            placeholder="Ej: Se llamó a emergencia, se aplicó hielo..."
                            value={form.action_taken}
                            onChange={(e) => setForm({ ...form, action_taken: e.target.value })}
                        />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                        <div className="grid gap-2">
                            <Label htmlFor="witnesses">Testigos (Opcional)</Label>
                            <Input
                                id="witnesses"
                                placeholder="Nombres..."
                                value={form.witnesses}
                                onChange={(e) => setForm({ ...form, witnesses: e.target.value })}
                            />
                        </div>
                        <div className="grid gap-2">
                            <Label htmlFor="injured">ID Usuario Lesionado (Opcional)</Label>
                            <Input
                                id="injured"
                                placeholder="UUID..."
                                value={form.injured_user_id}
                                onChange={(e) => setForm({ ...form, injured_user_id: e.target.value })}
                            />
                        </div>
                    </div>
                </div>
                <DialogFooter>
                    <Button type="submit" variant="destructive" onClick={handleSubmit} disabled={loading}>
                        {loading ? "Registrando..." : "Confirmar Reporte"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
