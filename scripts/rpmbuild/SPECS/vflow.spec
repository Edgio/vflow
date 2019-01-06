Name:       vflow
Version:    %VERSION%
Release:    0
Group:      Application
URL:        https://github.com/VerizonDigital/vflow
License:    Apache-2
Summary:    IPFIX/sFlow/Netflow collector
Source0:    vflow
Source1:    vflow_stress
Source2:    vflow.conf
Source3:    mq.conf
Source4:    vflow.service
Source5:    license
Source6:    notice
Source7:    vflow.logrotate

%description
High-performance, scalable and reliable IPFIX, sFlow and Netflow collector

%prep

%install
rm -rf %{buildroot}

mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/usr/local/vflow
mkdir -p %{buildroot}/usr/share/doc/vflow
mkdir -p %{buildroot}/etc/vflow
mkdir -p %{buildroot}/etc/init.d
mkdir -p %{buildroot}/etc/logrotate.d
cp -Rf %{SOURCE0} %{buildroot}/usr/bin/
cp -Rf %{SOURCE1} %{buildroot}/usr/bin/
cp -Rf %{SOURCE2} %{buildroot}/etc/vflow/
cp -Rf %{SOURCE3} %{buildroot}/etc/vflow/
cp -Rf %{SOURCE4} %{buildroot}/etc/init.d/vflow
cp -Rf %{SOURCE5} %{buildroot}/usr/share/doc/vflow/
cp -Rf %{SOURCE6} %{buildroot}/usr/share/doc/vflow/
cp -Rf %{SOURCE7} %{buildroot}/etc/logrotate.d/vflow

%files
/usr/bin/vflow
/usr/bin/vflow_stress
/etc/vflow/vflow.conf
/etc/vflow/mq.conf
/etc/init.d/vflow
/etc/logrotate.d/vflow
/usr/share/doc/vflow/*

%clean
rm -rf $RPM_BUILD_ROOT
