"use client";

import { useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
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
        password: "", // Added
        confirm_password: "", // Added
    });

    // GDPR Consent checkboxes
    const [acceptTerms, setAcceptTerms] = useState(false);
    const [acceptPrivacy, setAcceptPrivacy] = useState(false);
    const [parentalConsent, setParentalConsent] = useState(false);

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

        if (formData.password !== formData.confirm_password) {
            setStatus("error");
            setErrorMessage("Las contraseñas no coinciden.");
            return;
        }

        const passwordRegex = /^(?=.*[A-Z])(?=.*\d).{8,}$/;
        if (!passwordRegex.test(formData.password)) {
            setStatus("error");
            setErrorMessage("La contraseña debe tener al menos 8 caracteres, incluir una mayúscula y un número.");
            return;
        }

        // GDPR: Validate all consents are given
        if (!acceptTerms || !acceptPrivacy || !parentalConsent) {
            setStatus("error");
            setErrorMessage("Debe aceptar todos los términos y dar consentimiento parental para continuar.");
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
                },
                password: formData.password, // Added
                // GDPR Consent fields
                accept_terms: acceptTerms,
                privacy_policy_version: "2026-01",
                parental_consent: parentalConsent,
            };

            const res = await fetch(`http://localhost:8081/api/v1/users/public/register-dependent?club_id=${clubID}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload)
            });

            if (!res.ok) {
                const error = await res.json();
                throw new Error(error.error || "Failed to register");
            }

            setStatus("success");
            // Optional: Auto login logic could go here
        } catch (err: unknown) {
            setStatus("error");
            if (err instanceof Error) {
                // Mapear errores comunes a español
                const errorMap: Record<string, string> = {
                    'Failed to register': 'No pudimos completar el registro. Verificá los datos.',
                    'email already exists': 'Ese correo ya está registrado en el club.',
                    'invalid club': 'El club especificado no existe. Verificá el enlace.',
                };
                // Check for known error patterns
                const lowerMessage = err.message.toLowerCase();
                let translatedMessage = err.message;

                for (const [key, value] of Object.entries(errorMap)) {
                    if (lowerMessage.includes(key.toLowerCase())) {
                        translatedMessage = value;
                        break;
                    }
                }

                setErrorMessage(translatedMessage);
            } else {
                setErrorMessage("Ocurrió un error inesperado. Por favor, intentá de nuevo.");
            }
        }
    };

    if (status === "success") {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="text-green-600">¡Registro Exitoso!</CardTitle>
                    <CardDescription>
                        Los datos de {formData.child_name} {formData.child_surname} han sido registrados correctamente.
                        Ya puedes iniciar sesión con tu email y nueva contraseña.
                    </CardDescription>
                </CardHeader>
                <CardFooter className="flex gap-2">
                    <Button onClick={() => window.location.reload()} variant="outline" className="w-full">
                        Registrar otro
                    </Button>
                    <Link href="/login" className="w-full">
                        <Button className="w-full">Iniciar Sesión</Button>
                    </Link>
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
                                <Label htmlFor="password">Contraseña</Label>
                                <Input id="password" name="password" type="password" value={formData.password} onChange={handleChange} required placeholder="********" />
                            </div>
                            <div>
                                <Label htmlFor="confirm_password">Confirmar Contraseña</Label>
                                <Input id="confirm_password" name="confirm_password" type="password" value={formData.confirm_password} onChange={handleChange} required placeholder="********" />
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

                    {/* GDPR Consent Section */}
                    <div className="space-y-3 pt-4 border-t">
                        <h3 className="text-sm font-medium text-gray-500 uppercase tracking-wider">Consentimientos Requeridos</h3>

                        <div className="flex items-start space-x-3">
                            <Checkbox
                                id="accept_terms"
                                checked={acceptTerms}
                                onCheckedChange={(checked) => setAcceptTerms(checked as boolean)}
                                required
                            />
                            <Label htmlFor="accept_terms" className="text-sm leading-relaxed cursor-pointer">
                                He leído y acepto los{" "}
                                <Link href="/legal/terms" target="_blank" className="text-blue-600 underline hover:text-blue-800">
                                    Términos y Condiciones
                                </Link>{" "}
                                del club. <span className="text-red-500">*</span>
                            </Label>
                        </div>

                        <div className="flex items-start space-x-3">
                            <Checkbox
                                id="accept_privacy"
                                checked={acceptPrivacy}
                                onCheckedChange={(checked) => setAcceptPrivacy(checked as boolean)}
                                required
                            />
                            <Label htmlFor="accept_privacy" className="text-sm leading-relaxed cursor-pointer">
                                He leído y acepto la{" "}
                                <Link href="/legal/privacy" target="_blank" className="text-blue-600 underline hover:text-blue-800">
                                    Política de Privacidad
                                </Link>{" "}
                                y autorizo el tratamiento de datos personales conforme a la misma. <span className="text-red-500">*</span>
                            </Label>
                        </div>

                        <div className="flex items-start space-x-3 bg-amber-50 p-3 rounded-lg border border-amber-200">
                            <Checkbox
                                id="parental_consent"
                                checked={parentalConsent}
                                onCheckedChange={(checked) => setParentalConsent(checked as boolean)}
                                required
                            />
                            <Label htmlFor="parental_consent" className="text-sm leading-relaxed cursor-pointer">
                                <span className="font-medium text-amber-800">Consentimiento Parental:</span>{" "}
                                <span className="text-amber-700">
                                    Como padre/madre o tutor legal, autorizo el registro de mi hijo/a y el tratamiento de sus datos personales,
                                    incluyendo datos de salud (certificados médicos) necesarios para la práctica deportiva. <span className="text-red-500">*</span>
                                </span>
                            </Label>
                        </div>

                        <p className="text-xs text-gray-500 mt-2">
                            <span className="text-red-500">*</span> Campos obligatorios.
                            Puede ejercer sus derechos de acceso, rectificación, supresión y portabilidad contactando al club.
                        </p>
                    </div>

                    {errorMessage && <div className="text-red-500 text-sm font-medium bg-red-50 p-2 rounded">{errorMessage}</div>}

                    <Button
                        type="submit"
                        className="w-full bg-brand-600 hover:bg-brand-700"
                        disabled={status === "loading" || !acceptTerms || !acceptPrivacy || !parentalConsent}
                    >
                        {status === "loading" ? "Registrando..." : "Registrar Jugador"}
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

