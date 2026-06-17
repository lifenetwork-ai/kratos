local claims = std.extVar('claims');

std.prune({
  identity: {
    traits: {
      email: if std.objectHas(claims, 'email') then claims.email else null,
      phone_number: if std.objectHas(claims, 'phone_number') then claims.phone_number else null,
      tenant: if std.objectHas(claims, 'tenant_name') then claims.tenant_name else null,
      lang: 'en',
    },
  },
})
