%global debug_package %{nil}

Name:		gbalancer
Version:	0.5.4
Release:	1%{?dist}
Summary:	gbalancer orchestration tool

License:	ASL 2.0
URL:		https://github.com/coreos/%{name}/
Source0:	https://github.com/coreos/%{name}/archive/v%{version}/%{name}-v%{version}.tar.gz
Source1:	gbalancer.service
#Source2:	gbalancer.socket

BuildRequires:	golang
BuildRequires:	systemd
BuildRequires:	golang(github.com/coreos/go-systemd/activation) = 2-1.el7

Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd

%description
gbalancer orchestration tool

%prep
%setup -q

%build
mkdir -p gopath/src/github.com/coreos
./build

%install
install -D -p  bin/gbalancer %{buildroot}%{_bindir}/gbalancer
install -D -p  bin/gbalancerctl %{buildroot}%{_bindir}/gbalancerctl
install -D -p -m 0644 %{SOURCE1} %{buildroot}%{_unitdir}/%{name}.service
#install -D -p -m 0644 %{SOURCE2} %{buildroot}%{_unitdir}/%{name}.socket

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun %{name}.service

%files
%{_bindir}/gbalancer
%{_bindir}/gbalancerctl
%{_unitdir}/%{name}.service
#%{_unitdir}/%{name}.socket
%doc LICENSE README.md Documentation/architecture.md

%changelog
* Thu Jul 14 2014 Albert Zhang <zhgwenming@gmail.com> - 0.5.4-1
- initial version

