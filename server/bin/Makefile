name=verifier-server
bin=/usr/local/bin/${name}
init=/etc/init.d/${name}-d


install: ${bin} ${init}

${bin}:
	cp ./server-linux-amd64 ${bin}
	chmod +x ${bin}
	setcap 'cap_net_bind_service=+ep' ${bin}

${init}:
	cp ./init.sh ${init}
	chmod +x ${init}

clean:
	rm -f ${bin} ${init}

.PHONY: install clean