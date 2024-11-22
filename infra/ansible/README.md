### IT automation via Ansible

You need to install ansible before you can use
them.

```shell
pip3 install ansible
```

All ansible operation should be done under the
current directory.

You may also need to install sshpass for using ssh connection
with password. 
referece: [ install sshpass ](https://stackoverflow.com/questions/42835626/ansible-to-use-the-ssh-connection-type-with-passwords-you-must-install-the-s)

if you want to test your ansible connection,
you can try
```shell
ansible -i development all -m ping
```

For deploying an ansible-playbook,
```shell
ansible-playbook -i development hajime_tokenfate_deploy.yaml

```
### How to pass Extra Variables to Ansible Playbook?
```shell
--extra-vars "fruit=apple"
--extra-vars '{"fruit":"apple"}'
--extra-vars "@file.json"
--extra-vars "@file.yml"
```


Today we’re talking about Ansible extra variables. 
The easiest way to pass Pass Variables value to Ansible Playbook 
in the command line is using the extra variables parameter of the “ansible-playbook” command. 
This is very useful to combine your Ansible Playbook with some pre-existent automation or script. Let me clarify that this is specific for variables, 
there is another way to look for environment variables The command line parameter is the --extra-vars 
"variable=value" and allows you to pass some value from the terminal to the playbook. You could specify also the 
parameter in JSON format or include a JSON or a YAML file.
