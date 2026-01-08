'use client';

import { useState, useEffect, useCallback } from 'react';
import { userService, FamilyGroup, User } from '@/services/user-service';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Users, Plus, Loader2, UserPlus, Search, Crown } from 'lucide-react';

export default function FamilyGroupPage() {
    const [familyGroup, setFamilyGroup] = useState<FamilyGroup | null>(null);
    const [loading, setLoading] = useState(true);
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isAddMemberModalOpen, setIsAddMemberModalOpen] = useState(false);
    const [groupName, setGroupName] = useState('');
    const [userSearch, setUserSearch] = useState('');
    const [searchResults, setSearchResults] = useState<User[]>([]);
    const [submitting, setSubmitting] = useState(false);

    const loadFamilyGroup = useCallback(async () => {
        setLoading(true);
        try {
            const group = await userService.getMyFamilyGroup();
            setFamilyGroup(group);
        } catch (error) {
            console.error('Failed to load family group', error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadFamilyGroup();
    }, [loadFamilyGroup]);

    const handleCreateGroup = async () => {
        if (!groupName.trim()) return;
        setSubmitting(true);
        try {
            const group = await userService.createFamilyGroup(groupName);
            setFamilyGroup(group);
            setIsCreateModalOpen(false);
            setGroupName('');
            loadFamilyGroup(); // Reload to get members
        } catch (error) {
            console.error('Failed to create family group', error);
            alert('Error al crear el grupo familiar.');
        } finally {
            setSubmitting(false);
        }
    };

    const handleSearchUsers = async () => {
        if (userSearch.length < 2) return;
        try {
            const results = await userService.searchUsers(userSearch);
            // Filter out users already in the group
            const filteredResults = results.filter(
                u => !familyGroup?.members?.some(m => m.id === u.id)
            );
            setSearchResults(filteredResults);
        } catch (error) {
            console.error('Failed to search users', error);
        }
    };

    const handleAddMember = async (userId: string) => {
        if (!familyGroup) return;
        setSubmitting(true);
        try {
            await userService.addFamilyMember(familyGroup.id, userId);
            setIsAddMemberModalOpen(false);
            setSearchResults([]);
            setUserSearch('');
            loadFamilyGroup();
        } catch (error) {
            console.error('Failed to add member', error);
            alert('Error al agregar miembro.');
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center h-64">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    return (
        <div className="space-y-6 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Grupo Familiar</h1>
                    <p className="text-muted-foreground">Gestiona tu grupo familiar para facturación unificada.</p>
                </div>
            </div>

            {!familyGroup ? (
                <Card className="border-dashed">
                    <CardContent className="flex flex-col items-center justify-center py-12 text-center">
                        <Users className="h-12 w-12 text-muted-foreground mb-4" />
                        <h3 className="text-lg font-semibold mb-2">No tienes un grupo familiar</h3>
                        <p className="text-muted-foreground mb-6 max-w-sm">
                            Crea un grupo familiar para vincular cuentas de hijos o familiares y unificar la facturación.
                        </p>
                        <Dialog open={isCreateModalOpen} onOpenChange={setIsCreateModalOpen}>
                            <DialogTrigger asChild>
                                <Button>
                                    <Plus className="mr-2 h-4 w-4" />
                                    Crear Grupo Familiar
                                </Button>
                            </DialogTrigger>
                            <DialogContent>
                                <DialogHeader>
                                    <DialogTitle>Crear Grupo Familiar</DialogTitle>
                                    <DialogDescription>
                                        Ingresa un nombre para identificar tu grupo familiar.
                                    </DialogDescription>
                                </DialogHeader>
                                <div className="py-4">
                                    <Label htmlFor="groupName">Nombre del Grupo</Label>
                                    <Input
                                        id="groupName"
                                        placeholder="ej. Familia García"
                                        value={groupName}
                                        onChange={(e) => setGroupName(e.target.value)}
                                    />
                                </div>
                                <DialogFooter>
                                    <Button variant="outline" onClick={() => setIsCreateModalOpen(false)}>Cancelar</Button>
                                    <Button onClick={handleCreateGroup} disabled={submitting || !groupName.trim()}>
                                        {submitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                        Crear
                                    </Button>
                                </DialogFooter>
                            </DialogContent>
                        </Dialog>
                    </CardContent>
                </Card>
            ) : (
                <>
                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between">
                            <div>
                                <CardTitle className="flex items-center gap-2">
                                    <Users className="h-5 w-5" />
                                    {familyGroup.name}
                                </CardTitle>
                                <CardDescription>
                                    {familyGroup.members?.length || 0} miembro(s) en el grupo
                                </CardDescription>
                            </div>
                            <Dialog open={isAddMemberModalOpen} onOpenChange={setIsAddMemberModalOpen}>
                                <DialogTrigger asChild>
                                    <Button size="sm">
                                        <UserPlus className="mr-2 h-4 w-4" />
                                        Agregar Miembro
                                    </Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>Agregar Miembro</DialogTitle>
                                        <DialogDescription>
                                            Busca un usuario existente para agregar a tu grupo familiar.
                                        </DialogDescription>
                                    </DialogHeader>
                                    <div className="py-4 space-y-4">
                                        <div className="flex gap-2">
                                            <Input
                                                placeholder="Buscar por nombre o email..."
                                                value={userSearch}
                                                onChange={(e) => setUserSearch(e.target.value)}
                                                onKeyDown={(e) => e.key === 'Enter' && handleSearchUsers()}
                                            />
                                            <Button variant="outline" size="icon" onClick={handleSearchUsers}>
                                                <Search className="h-4 w-4" />
                                            </Button>
                                        </div>
                                        {searchResults.length > 0 && (
                                            <div className="border rounded-md max-h-48 overflow-y-auto">
                                                {searchResults.map((user) => (
                                                    <div
                                                        key={user.id}
                                                        className="p-3 flex justify-between items-center border-b last:border-b-0 hover:bg-muted"
                                                    >
                                                        <div>
                                                            <div className="font-medium">{user.name}</div>
                                                            <div className="text-sm text-muted-foreground">{user.email}</div>
                                                        </div>
                                                        <Button
                                                            size="sm"
                                                            onClick={() => handleAddMember(user.id)}
                                                            disabled={submitting}
                                                        >
                                                            Agregar
                                                        </Button>
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                    <DialogFooter>
                                        <Button variant="outline" onClick={() => setIsAddMemberModalOpen(false)}>Cerrar</Button>
                                    </DialogFooter>
                                </DialogContent>
                            </Dialog>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-3">
                                {familyGroup.members?.map((member) => (
                                    <div
                                        key={member.id}
                                        className="flex items-center justify-between p-3 rounded-lg border"
                                    >
                                        <div className="flex items-center gap-3">
                                            <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center">
                                                {member.name.charAt(0).toUpperCase()}
                                            </div>
                                            <div>
                                                <div className="font-medium flex items-center gap-2">
                                                    {member.name}
                                                    {member.id === familyGroup.head_user_id && (
                                                        <Badge variant="secondary" className="text-xs">
                                                            <Crown className="h-3 w-3 mr-1" />
                                                            Titular
                                                        </Badge>
                                                    )}
                                                </div>
                                                <div className="text-sm text-muted-foreground">{member.email}</div>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </>
            )}
        </div>
    );
}
