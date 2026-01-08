"use client";

import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export default function PrivacyPage() {
    return (
        <div className="min-h-screen bg-gray-50 py-12 px-4">
            <div className="max-w-4xl mx-auto">
                <Card>
                    <CardHeader>
                        <CardTitle className="text-2xl">Política de Privacidad</CardTitle>
                        <p className="text-sm text-gray-500">Versión: 2026-01 | Última actualización: Enero 2026</p>
                    </CardHeader>
                    <CardContent className="prose prose-gray max-w-none">
                        <h2>1. Responsable del Tratamiento</h2>
                        <p>
                            El responsable del tratamiento de sus datos personales es el club deportivo al que
                            se ha registrado. Puede contactarnos a través de la administración del club.
                        </p>

                        <h2>2. Datos que Recopilamos</h2>
                        <p>Recopilamos los siguientes tipos de datos personales:</p>
                        <ul>
                            <li><strong>Datos de identificación:</strong> Nombre, email, teléfono, fecha de nacimiento</li>
                            <li><strong>Datos de contacto de emergencia:</strong> Nombre y teléfono de contacto</li>
                            <li><strong>Datos de salud:</strong> Certificados médicos, información de seguro</li>
                            <li><strong>Datos de menores:</strong> Con consentimiento parental verificado</li>
                            <li><strong>Datos de uso:</strong> Reservas, asistencia, pagos</li>
                        </ul>

                        <h2>3. Finalidad del Tratamiento</h2>
                        <p>Sus datos personales son tratados para:</p>
                        <ul>
                            <li>Gestionar su membresía y acceso a instalaciones</li>
                            <li>Procesar reservas y pagos</li>
                            <li>Garantizar la seguridad deportiva (certificados médicos)</li>
                            <li>Organizar competencias y eventos</li>
                            <li>Comunicaciones relacionadas con el servicio</li>
                        </ul>

                        <h2>4. Base Legal</h2>
                        <p>El tratamiento de sus datos se basa en:</p>
                        <ul>
                            <li><strong>Consentimiento:</strong> Otorgado al aceptar estos términos</li>
                            <li><strong>Ejecución contractual:</strong> Prestación de servicios deportivos</li>
                            <li><strong>Obligación legal:</strong> Requisitos de seguridad deportiva</li>
                        </ul>

                        <h2>5. Datos de Menores (GDPR Art. 8)</h2>
                        <p>
                            El tratamiento de datos de menores de 18 años requiere consentimiento parental
                            verificable. Los padres/tutores son responsables de autorizar el registro.
                        </p>

                        <h2>6. Datos de Salud (GDPR Art. 9)</h2>
                        <p>
                            Los certificados médicos son tratados como categoría especial de datos.
                            El acceso está restringido al personal médico autorizado y solo se utiliza
                            para verificar aptitud deportiva.
                        </p>

                        <h2>7. Sus Derechos</h2>
                        <p>Conforme al GDPR/LGPD, usted tiene derecho a:</p>
                        <ul>
                            <li><strong>Acceso:</strong> Solicitar una copia de sus datos</li>
                            <li><strong>Rectificación:</strong> Corregir datos inexactos</li>
                            <li><strong>Supresión:</strong> Solicitar eliminación de sus datos (Art. 17)</li>
                            <li><strong>Portabilidad:</strong> Recibir sus datos en formato portable (Art. 20)</li>
                            <li><strong>Oposición:</strong> Oponerse al tratamiento</li>
                            <li><strong>Limitación:</strong> Limitar el tratamiento en ciertos casos</li>
                        </ul>
                        <p>
                            Para ejercer estos derechos, contacte a la administración del club o utilice
                            las opciones disponibles en su perfil de usuario.
                        </p>

                        <h2>8. Retención de Datos</h2>
                        <p>
                            Conservamos sus datos mientras mantenga una membresía activa y durante el período
                            requerido por obligaciones legales. Tras la baja, los datos son anonimizados.
                        </p>

                        <h2>9. Seguridad</h2>
                        <p>
                            Implementamos medidas técnicas y organizativas para proteger sus datos, incluyendo
                            cifrado, control de acceso y auditoría de accesos a datos sensibles.
                        </p>

                        <h2>10. Contacto</h2>
                        <p>
                            Para consultas sobre privacidad o para ejercer sus derechos, contacte a la
                            administración del club.
                        </p>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
