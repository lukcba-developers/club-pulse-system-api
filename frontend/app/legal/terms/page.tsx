"use client";

import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export default function TermsPage() {
    return (
        <div className="min-h-screen bg-gray-50 py-12 px-4">
            <div className="max-w-4xl mx-auto">
                <Card>
                    <CardHeader>
                        <CardTitle className="text-2xl">Términos y Condiciones</CardTitle>
                        <p className="text-sm text-gray-500">Última actualización: Enero 2026</p>
                    </CardHeader>
                    <CardContent className="prose prose-gray max-w-none">
                        <h2>1. Aceptación de los Términos</h2>
                        <p>
                            Al registrarse y utilizar los servicios del club, usted acepta estar vinculado por estos
                            Términos y Condiciones. Si no está de acuerdo con alguna parte de estos términos,
                            no debe utilizar nuestros servicios.
                        </p>

                        <h2>2. Descripción del Servicio</h2>
                        <p>
                            El club ofrece servicios deportivos que incluyen:
                        </p>
                        <ul>
                            <li>Acceso a instalaciones deportivas</li>
                            <li>Participación en actividades y entrenamientos</li>
                            <li>Inscripción en torneos y competencias</li>
                            <li>Alquiler de equipamiento deportivo</li>
                        </ul>

                        <h2>3. Registro y Cuenta</h2>
                        <p>
                            Para utilizar nuestros servicios, debe proporcionar información precisa y actualizada.
                            Es responsable de mantener la confidencialidad de su cuenta y contraseña.
                        </p>

                        <h2>4. Menores de Edad</h2>
                        <p>
                            Los menores de 18 años requieren autorización de un padre o tutor legal para registrarse.
                            El padre/tutor es responsable de las actividades del menor en el club.
                        </p>

                        <h2>5. Pagos y Membresías</h2>
                        <p>
                            Las cuotas de membresía y otros pagos deben realizarse según los plazos establecidos.
                            El incumplimiento puede resultar en la suspensión de servicios.
                        </p>

                        <h2>6. Responsabilidad</h2>
                        <p>
                            El club no se hace responsable de lesiones o daños que ocurran durante la práctica
                            deportiva, salvo en casos de negligencia comprobada por parte del club.
                        </p>

                        <h2>7. Modificaciones</h2>
                        <p>
                            El club se reserva el derecho de modificar estos términos en cualquier momento.
                            Los cambios serán notificados y entrarán en vigor inmediatamente.
                        </p>

                        <h2>8. Contacto</h2>
                        <p>
                            Para consultas sobre estos términos, contacte a la administración del club.
                        </p>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
