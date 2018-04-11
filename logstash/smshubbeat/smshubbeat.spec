Name:		smshubbeat
Version:    %{version}
Release:    %{buildnumber}
Summary:    smshub 3 redis kpi logstash input beat	
Group:		Applications/System
License:	TATA
Source:     smshubbeat-%{version}.tgz	
Packager:   Daniel Lavalliere daniel.lavalliere@tatacommunications.com
BuildRoot:  /var/tmp/%{name}-buildroot


%description


%prep
%setup -q


%build
export GOPATH=%{githome}
cd %{buildpath}
%{__make}


%install
cd %{homesrc}
rm -rf $RPM_BUILD_ROOT
%{__install} -m 755 -d %{buildroot}/etc/init.d
%{__install} -m 755 -d %{buildroot}/etc/smshubbeat
%{__install} -m 755 -d %{buildroot}/lib/systemd/system
%{__install} -m 755 -d %{buildroot}/usr/bin
%{__install} -m 755 -d %{buildroot}/usr/share/smshubbeat
%{__install} -m 755 -d %{buildroot}/usr/share/smshubbeat/bin
%{__install} -m 755 smshubbeat.init %{buildroot}/etc/init.d/smshubbeat
%{__install} -m 600 smshubbeat.yml %{buildroot}/etc/smshubbeat/smshubbeat.yml
%{__install} -m 644 smshubbeat.full.yml %{buildroot}/etc/smshubbeat/smshubbeat.full.yml
%{__install} -m 644 smshubbeat.template-es2x.json %{buildroot}/etc/smshubbeat/smshubbeat.template-es2x.json
%{__install} -m 644 smshubbeat.template.json %{buildroot}/etc/smshubbeat/smshubbeat.template.json
%{__install} -m 600 smshubbeat.lua %{buildroot}/etc/smshubbeat/smshubbeat.lua
%{__install} -m 644 smshubbeat.service %{buildroot}/lib/systemd/system/smshubbeat.service
%{__install} -m 755 smshubbeat.sh %{buildroot}/usr/bin/smshubbeat.sh
%{__install} -m 644 NOTICE %{buildroot}/usr/share/smshubbeat/NOTICE
%{__install} -m 644 README.md.rpm %{buildroot}/usr/share/smshubbeat/README.md
cd %{buildpath}
%{__install} -m 755 smshubbeat %{buildroot}/usr/share/smshubbeat/bin/smshubbeat


%clean
rm -rf $RPM_BUILD_ROOT


%files
%defattr(-,root,root)
%attr(-, root, root) /etc/init.d/smshubbeat
%attr(-, root, root) /usr/share/smshubbeat/bin/smshubbeat
%attr(-, root, root) /lib/systemd/system/smshubbeat.service
%attr(-, root, root) /usr/bin/smshubbeat.sh
%attr(-, root, root) /usr/share/smshubbeat/NOTICE
%attr(-, root, root) /usr/share/smshubbeat/README.md
%config(noreplace) /etc/smshubbeat/smshubbeat.yml
%config(noreplace) /etc/smshubbeat/smshubbeat.full.yml
%config(noreplace) /etc/smshubbeat/smshubbeat.lua
%config(noreplace) /etc/smshubbeat/smshubbeat.template-es2x.json
%config(noreplace) /etc/smshubbeat/smshubbeat.template.json


