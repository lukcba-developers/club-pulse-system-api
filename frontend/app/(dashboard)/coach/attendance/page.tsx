'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/hooks/use-auth';
import { useRouter } from 'next/navigation';
import { attendanceService, AttendanceList, AttendanceRecord } from '@/services/attendance-service';
import { disciplineService, Discipline, TrainingGroup } from '@/services/discipline-service';
import { Check, X, Clock, Users, AlertTriangle } from 'lucide-react';


export default function CoachAttendancePage() {
    const { user, loading } = useAuth();
    const router = useRouter();
    const [disciplines, setDisciplines] = useState<Discipline[]>([]);
    const [selectedDiscipline, setSelectedDiscipline] = useState<string>('');
    const [groups, setGroups] = useState<TrainingGroup[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState<string>('');

    const [attendanceList, setAttendanceList] = useState<AttendanceList | null>(null);
    const [loadingList, setLoadingList] = useState(false);
    const [savingId, setSavingId] = useState<string | null>(null);

    useEffect(() => {
        if (!loading && !user) {
            router.push('/login');
        }
        if (user) {
            fetchDisciplines();
        }
    }, [user, loading, router]);

    const fetchDisciplines = async () => {
        try {
            const data = await disciplineService.listDisciplines();
            setDisciplines(data);
        } catch (error) {
            console.error('Failed to fetch disciplines', error);
        }
    };

    const fetchGroups = async (disciplineId: string) => {
        try {
            const data = await disciplineService.listGroups(disciplineId);
            setGroups(data);
        } catch (error) {
            console.error('Failed to fetch groups', error);
        }
    };

    const fetchList = async (groupId: string) => {
        const group = groups.find(g => g.id === groupId);
        if (!group) return;

        setLoadingList(true);
        try {
            const today = new Date().toISOString().split('T')[0];
            const data = await attendanceService.getTrainingGroupList(
                group.id,
                group.name,
                group.category,
                today
            );
            setAttendanceList(data);
        } catch (error) {
            console.error('Failed to fetch attendance list', error);
            setAttendanceList(null);
        } finally {
            setLoadingList(false);
        }
    };

    const handleDisciplineChange = (id: string) => {
        setSelectedDiscipline(id);
        setSelectedGroupId('');
        setAttendanceList(null);
        if (id) fetchGroups(id);
        else setGroups([]);
    };

    const handleGroupChange = (id: string) => {
        setSelectedGroupId(id);
        if (id) fetchList(id);
        else setAttendanceList(null);
    };

    const handleMark = async (record: AttendanceRecord, newStatus: string) => {
        if (!attendanceList) return;
        setSavingId(record.user_id);
        try {
            await attendanceService.markAttendance(attendanceList.id, record.user_id, newStatus);
            // Update local state
            setAttendanceList(prev => {
                if (!prev) return null;
                return {
                    ...prev,
                    records: prev.records.map(r =>
                        r.user_id === record.user_id ? { ...r, status: newStatus as AttendanceRecord['status'] } : r
                    )
                };
            });
        } catch (error) {
            console.error('Failed to mark attendance', error);
        } finally {
            setSavingId(null);
        }
    };

    if (loading || !user) {
        return (
            <div className="h-full flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-600"></div>
            </div>
        );
    }

    return (
        <div className="max-w-4xl mx-auto">
            <div className="mb-8">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
                    <Users className="h-6 w-6" />
                    Control de Asistencia
                </h1>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    Selecciona una disciplina y grupo para gestionar la asistencia.
                </p>
            </div>

            {/* Selectors */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        Disciplina
                    </label>
                    <select
                        value={selectedDiscipline}
                        onChange={(e) => handleDisciplineChange(e.target.value)}
                        className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-brand-500"
                    >
                        <option value="">-- Seleccionar Disciplina --</option>
                        {disciplines.map(d => (
                            <option key={d.id} value={d.id}>{d.name}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        Grupo de Entrenamiento
                    </label>
                    <select
                        value={selectedGroupId}
                        onChange={(e) => handleGroupChange(e.target.value)}
                        disabled={!selectedDiscipline}
                        className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-brand-500 disabled:opacity-50"
                    >
                        <option value="">-- Seleccionar Grupo --</option>
                        {groups.map(g => (
                            <option key={g.id} value={g.id}>{g.name} ({g.category})</option>
                        ))}
                    </select>
                </div>
            </div>

            {/* Attendance List */}
            {loadingList && (
                <div className="flex justify-center py-12">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-600"></div>
                </div>
            )}

            {!loadingList && attendanceList && (
                <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
                    <div className="px-6 py-4 border-b border-gray-100 dark:border-gray-700">
                        <h2 className="font-semibold text-gray-900 dark:text-white">
                            Lista del {new Date(attendanceList.date).toLocaleDateString('es-AR')} - {attendanceList.group}
                        </h2>
                        <p className="text-sm text-gray-500">{attendanceList.records?.length || 0} alumnos registrados</p>
                    </div>

                    {(!attendanceList.records || attendanceList.records.length === 0) ? (
                        <div className="p-6 text-center text-gray-500">
                            No hay alumnos en esta categoría.
                        </div>
                    ) : (
                        <ul className="divide-y divide-gray-100 dark:divide-gray-700">
                            {attendanceList.records.map((record) => (
                                <li key={record.user_id} className="px-6 py-4 flex items-center justify-between gap-4">
                                    <div>
                                        <p className="font-medium text-gray-900 dark:text-white">
                                            {record.user?.name || record.user_id}
                                            {record.has_debt && (
                                                <span title="Deuda Pendiente">
                                                    <AlertTriangle className="h-4 w-4 text-red-500 inline ml-2" />
                                                </span>
                                            )}
                                        </p>
                                        <span className={`text-xs px-2 py-0.5 rounded-full ${record.status === 'PRESENT' ? 'bg-green-100 text-green-700' :
                                            record.status === 'LATE' ? 'bg-yellow-100 text-yellow-700' :
                                                'bg-red-100 text-red-700'
                                            }`}>
                                            {record.status === 'PRESENT' ? 'Presente' :
                                                record.status === 'LATE' ? 'Tarde' : 'Ausente'}
                                        </span>
                                    </div>
                                    <div className="flex gap-2">
                                        <button
                                            onClick={() => handleMark(record, 'PRESENT')}
                                            disabled={savingId === record.user_id}
                                            className={`p-2 rounded-lg transition ${record.status === 'PRESENT' ? 'bg-green-500 text-white' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-green-100'}`}
                                            title="Presente"
                                        >
                                            <Check className="h-5 w-5" />
                                        </button>
                                        <button
                                            onClick={() => handleMark(record, 'LATE')}
                                            disabled={savingId === record.user_id}
                                            className={`p-2 rounded-lg transition ${record.status === 'LATE' ? 'bg-yellow-500 text-white' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-yellow-100'}`}
                                            title="Tarde"
                                        >
                                            <Clock className="h-5 w-5" />
                                        </button>
                                        <button
                                            onClick={() => handleMark(record, 'ABSENT')}
                                            disabled={savingId === record.user_id}
                                            className={`p-2 rounded-lg transition ${record.status === 'ABSENT' ? 'bg-red-500 text-white' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-red-100'}`}
                                            title="Ausente"
                                        >
                                            <X className="h-5 w-5" />
                                        </button>
                                    </div>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>
            )}

            {!loadingList && selectedGroupId && !attendanceList && (
                <div className="text-center py-12 text-gray-500">
                    No se pudo cargar la lista. Intenta de nuevo.
                </div>
            )}
        </div>
    );
}
