'use client';

import { useState, useEffect, useCallback } from 'react';
import { userService, User } from '@/services/user-service';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Search } from 'lucide-react';

export default function UsersPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSearch = useCallback(async () => { // Wrapped in useCallback
        setLoading(true);
        try {
            const results = await userService.searchUsers(searchQuery);
            // Handle { data: [...] } vs [...]
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const usersData = (results as any).data || results || [];
            if (Array.isArray(usersData)) {
                setUsers(usersData);
            } else {
                setUsers([]);
            }
        } catch (error) {
            console.error('Failed to fetch users:', error);
        } finally {
            setLoading(false);
        }
    }, [searchQuery]); // Dependency: searchQuery

    useEffect(() => {
        // Initial load
        handleSearch();
    }, [handleSearch]); // Added dependency

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-bold tracking-tight">Usuarios</h1>
            </div>

            <div className="flex gap-2">
                <Input
                    placeholder="Buscar por nombre o email..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                    className="max-w-sm"
                />
                <Button onClick={handleSearch} disabled={loading}>
                    <Search className="mr-2 h-4 w-4" />
                    Buscar
                </Button>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {users.map((user) => (
                    <Card key={user.id}>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-lg font-medium">{user.name}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-sm text-muted-foreground">{user.email}</div>
                            <div className="text-xs mt-2  bg-secondary inline-block px-2 py-1 rounded-full">{user.role}</div>
                        </CardContent>
                    </Card>
                ))}
                {!loading && users.length === 0 && (
                    <div className="col-span-full text-center text-muted-foreground py-10">
                        No se encontraron usuarios.
                    </div>
                )}
            </div>
        </div>
    );
}
