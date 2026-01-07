"use client";

import { useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/card";

function RegisterForm() {
    const searchParams = useSearchParams();
    const clubIdParam = searchParams.get("club_id");
    const [clubID, setClubID] = useState(clubIdParam || "");

    const [formData, setFormData] = useState({
        parent_email: "",
        parent_name: "",
        parent_phone: "",
        child_name: "",
        child_surname: "",
        child_dob: "",
        sport: "TENNIS",
    });

    const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
    const [errorMessage, setErrorMessage] = useState("");

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setStatus("loading");
        setErrorMessage("");

        if (!clubID) {
            setStatus("error");
            setErrorMessage("Se requiere el ID del club. Verifique el enlace.");
            return;
        }

        try {
            const payload = {
                parent_email: formData.parent_email,
                parent_name: formData.parent_name,
                parent_phone: formData.parent_phone,
                child_name: formData.child_name,
                child_surname: formData.child_surname,
                child_dob: new Date(formData.child_dob).toISOString(),
                sports_preferences: {
                    primary: formData.sport
                }
            };

            const res = await fetch(`http://localhost:8080/api/v1/users/public/register-dependent?club_id=${clubID}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload)
            });

            if (!res.ok) {
                const error = await res.json();
                throw new Error(error.error || "Failed to register");
            }

            setStatus("success");
        } catch (err: any) {
            setStatus("error");
            setErrorMessage(err.message);
        }
    };

    if (status === "success") {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="text-green-600">¡Registro Exitoso!</CardTitle>
                    <CardDescription>
                        Los datos de {formData.child_name} {formData.child_surname} han sido registrados correctamente.
                    </CardDescription>
                </CardHeader>
                <CardFooter>
                    <Button onClick={() => window.location.reload()} variant="outline" className="w-full">
                        Registrar otro
                    </Button>
                </CardFooter>
            </Card>
        )
    }

    return (
        <Card className="w-full">
            <CardHeader>
                <CardTitle>Registro de Jugador</CardTitle>
                <CardDescription>
                    Completa el formulario para registrar a tu hijo/a en el club.
                </CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    {!clubIdParam && (
                        <div className="space-y-2">
                            <Label htmlFor="club_id">ID del Club (Requerido)</Label>
                            <Input id="club_id" name="club_id" value={clubID} onChange={(e) => setClubID(e.target.value)} required placeholder="UUID" />
                        </div>
                    )}

                    <div className="space-y-2">
                        <h3 className="text-sm font-medium text-gray-500 uppercase tracking-wider">Datos del Responsable</h3>
                        <div className="grid grid-cols-1 gap-3">
                            <div>
                                <Label htmlFor="parent_email">Email</Label>
                                <Input id="parent_email" name="parent_email" type="email" value={formData.parent_email} onChange={handleChange} required placeholder="tu@email.com" />
                            </div>
                            <div>
                                <Label htmlFor="parent_name">Nombre Completo</Label>
                                <Input id="parent_name" name="parent_name" value={formData.parent_name} onChange={handleChange} required placeholder="Juan Pérez" />
                            </div>
                            <div>
                                <Label htmlFor="parent_phone">Teléfono (WhatsApp)</Label>
                                <Input id="parent_phone" name="parent_phone" type="tel" value={formData.parent_phone} onChange={handleChange} required placeholder="+54 9 11 ..." />
                            </div>
                        </div>
                    </div>

                    <div className="space-y-2 pt-2 border-t">
                        <h3 className="text-sm font-medium text-gray-500 uppercase tracking-wider">Datos del Jugador</h3>
                        <div className="grid grid-cols-1 gap-3">
                            <div className="grid grid-cols-2 gap-2">
                                <div>
                                    <Label htmlFor="child_name">Nombre</Label>
                                    <Input id="child_name" name="child_name" value={formData.child_name} onChange={handleChange} required placeholder="Mateo" />
                                </div>
                                <div>
                                    <Label htmlFor="child_surname">Apellido</Label>
                                    <Input id="child_surname" name="child_surname" value={formData.child_surname} onChange={handleChange} required placeholder="Pérez" />
                                </div>
                            </div>
                            <div>
                                <Label htmlFor="child_dob">Fecha de Nacimiento</Label>
                                <Input id="child_dob" name="child_dob" type="date" value={formData.child_dob} onChange={handleChange} required />
                            </div>
                            <div>
                                <Label htmlFor="sport">Deporte de Interés</Label>
                                <select
                                    id="sport"
                                    name="sport"
                                    className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    value={formData.sport}
                                    onChange={handleChange}
                                >
                                    <option value="TENNIS">Tenis</option>
                                    <option value="FOOTBALL">Fútbol</option>
                                    <option value="PADEL">Pádel</option>
                                    <option value="HOCKEY">Hockey</option>
                                </select>
                            </div>
                        </div>
                    </div>

                    {errorMessage && <div className="text-red-500 text-sm font-medium bg-red-50 p-2 rounded">{errorMessage}</div>}

                    <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700" disabled={status === "loading"}>
                        {status === "loading" ? "Registrando..." : "Registrar"}
                    </Button>
                </form>
            </CardContent>
        </Card>
    );
}

export default function RegisterPlayerPage() {
    return (
        <Suspense fallback={<div className="text-center p-4">Cargando...</div>}>
            <RegisterForm />
        </Suspense>
    )
}
