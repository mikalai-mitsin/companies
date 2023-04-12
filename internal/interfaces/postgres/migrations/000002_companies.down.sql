DROP TABLE public.companies;

DELETE
FROM public.permissions
WHERE id IN (
    'company_list',
    'company_detail',
    'company_create',
    'company_update',
    'company_delete'
);
