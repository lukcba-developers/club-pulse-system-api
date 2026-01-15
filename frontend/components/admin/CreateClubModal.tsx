'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { clubService } from '@/services/club-service';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

const clubSchema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    slug: z.string().min(2).regex(/^[a-z0-9-]+$/, 'Slug must be lowercase alphanumeric with hyphens'),
    domain: z.string().optional(),
    logo_url: z.string().url('Invalid URL').optional().or(z.literal('')),
    theme_config: z.string().refine((val) => {
        if (!val) return true;
        try { JSON.parse(val); return true; } catch { return false; }
    }, 'Invalid JSON format').optional(),
    settings: z.string().refine((val) => {
        if (!val) return true;
        try { JSON.parse(val); return true; } catch { return false; }
    }, 'Invalid JSON format').optional(),
});

type ClubFormData = z.infer<typeof clubSchema>;

interface CreateClubModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onSuccess: () => void;
}

export function CreateClubModal({ open, onOpenChange, onSuccess }: CreateClubModalProps) {
    const [loading, setLoading] = useState(false);
    const { register, handleSubmit, formState: { errors }, reset } = useForm<ClubFormData>({
        resolver: zodResolver(clubSchema),
        defaultValues: {
            name: '',
            slug: '',
            domain: '',
            logo_url: '',
            theme_config: '',
            settings: ''
        }
    });

    const onSubmit = async (data: ClubFormData) => {
        setLoading(true);
        try {
            await clubService.createClub({
                ...data,
                domain: data.domain || undefined,
                logo_url: data.logo_url || undefined,
                theme_config: data.theme_config || undefined,
                settings: data.settings || undefined,
            });
            onSuccess();
            onOpenChange(false);
            reset();
        } catch (err: unknown) {
            console.error(err);
            const errorMessage = err instanceof Error ? err.message : 'Failed to create club';
            // handle global error if needed, but hook form handles field errors
            // For now just alert or log, as simple modal
            alert(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Create New Club</DialogTitle>
                    <DialogDescription>
                        Fill in the details to create a new club.
                    </DialogDescription>
                </DialogHeader>
                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="grid gap-2">
                        <Label htmlFor="name">Name</Label>
                        <Input id="name" {...register('name')} />
                        {errors.name && <p className="text-red-500 text-xs">{errors.name.message}</p>}
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="slug">Slug (URL Identifier)</Label>
                        <Input id="slug" {...register('slug')} />
                        {errors.slug && <p className="text-red-500 text-xs">{errors.slug.message}</p>}
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="domain">Custom Domain</Label>
                        <Input id="domain" {...register('domain')} placeholder="club.com" />
                        {errors.domain && <p className="text-red-500 text-xs">{errors.domain.message}</p>}
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="logo_url">Logo URL</Label>
                        <Input id="logo_url" {...register('logo_url')} placeholder="https://..." />
                        {errors.logo_url && <p className="text-red-500 text-xs">{errors.logo_url.message}</p>}
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="theme_config">Theme Config (JSON)</Label>
                        <Textarea id="theme_config" {...register('theme_config')} placeholder='{"primary": "#ff0000"}' />
                        {errors.theme_config && <p className="text-red-500 text-xs">{errors.theme_config.message}</p>}
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="settings">Settings (JSON)</Label>
                        <Textarea id="settings" {...register('settings')} placeholder='{"timezone": "UTC"}' />
                        {errors.settings && <p className="text-red-500 text-xs">{errors.settings.message}</p>}
                    </div>
                    <DialogFooter>
                        <Button type="submit" disabled={loading}>
                            {loading ? 'Creating...' : 'Create Club'}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent >
        </Dialog >
    );
}
