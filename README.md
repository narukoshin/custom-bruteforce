<img src="https://c.tenor.com/gOP4dRPvzWcAAAAi/angry-mafumafu.gif" align="right" width="160">
<h1>🧪 A little bit different brute-force tool</h1>
<div>
  <img src="https://img.shields.io/github/go-mod/go-version/narukoshin/custom-bruteforce">
  <img src="https://img.shields.io/github/v/release/narukoshin/custom-bruteforce">
  <img src="https://img.shields.io/github/last-commit/narukoshin/custom-bruteforce">
  <img src="https://img.shields.io/github/contributors/narukoshin/custom-bruteforce">
  <br><br>
  <div>
    <a target="_blank" href="https://twitter.com/narukoshin"><img src="https://media4.giphy.com/media/iFUiSYMNPvIJZDpMKN/giphy.gif?cid=ecf05e471v5jn6vuhczu1tflu2wm7qt11atwybfwcgaqxz38&rid=giphy.gif&ct=s" width="120"></a>
    <a target="_blank" href="https://instagram.com/naru.koshin"><img src="https://media1.giphy.com/media/Wu9Graz2W46frtHFKc/giphy.gif?cid=ecf05e47h46mbuhq40rgevni5rbxgadpw5icrr71vr9nu8d4&rid=giphy.gif&ct=s" width="120"></a>
  </div>
</div>
<h1>📚 Getting started</h1>
<p>To download this tool, type the command below:</p>

```sh
  git clone https://github.com/narukoshin/custom-bruteforce
```
... or download binaries from the releases page.

<h1>📅 TODO</h1>

- [ ] Proxy Feature
- [x] ~~Idea about `import` option where you can import config file with the custom name like `import: my_website.yml`~~ <br>
      + added in commit: <a href="https://github.com/narukoshin/custom-bruteforce/commit/823b14f907ce92a44d69174510f681ba0da31c6e">`823b14f`</a><br>
      + changelog: <a href="https://github.com/narukoshin/custom-bruteforce/releases/tag/v2.3-beta">`v2.3-beta`</a>
- [ ] Email notifications

💭 If you have any suggestion about new features, please open a new issue with the enhancement label.

<h1>⚙ Creating configuration</h1>
<p>Before you start using the tool, you need to create a config file called <code>config.yml</code></p>

```sh
touch config.yml
... or
vim config.yml
... or you can use pre-made config
mv config.sample.yml config.yml
```

<p>Next, you need to fill the config file with the information about your target to brute-force.</p>

```yaml
config.yml

site:
    host: https://website.com/login # the login page that you want to crack.
    method: POST # request method for making a request
bruteforce:
    field: password # the field that you want to brute-force (important)
    
    # there is 3 ways from where you can load a wordlist
    # method 1 - from the file
    from: file
    file: /usr/share/wordlists/rockyou.txt # the path, where is your wordlist located at
    # method 2 - from the list
    from: list
    list:
        - password1
        - password2
        - password3
    # method 3 - from the stdin
    from: stdin
```
<p>When you are using <b>stdin</b> method, type the command as shown below:</p>

```sh
    # example 1
    cat /usr/share/wordlists/rockyou.txt | ./linux
    # example 2
    crunch 8 8 0123456789 | ./linux
    # ...etc
```

```yaml
config.yml

    # Next, you need to specify how many threads you want to use. The default value is 5
    threads: 30
    
    # If you don't want to see messages like "trying password...", you can turn it off with the option below:
    # It's optional, so if you don't need to turn it off, you can skip this option
    no_verbose: true
    
    # By default, when the tool finds out the password, the password will be printed on the screen, 
    # ...but if you want you can set it to save in the file
    output: /home/naru/my_target/password.txt
    
# Setting the headers
# There's nothing difficult
headers:
    - name: Content-Type
      value: application/x-www-form-urlencoded; charset=utf-8
    - name: User-Agent
      value: Mozilla/5.0 (X11; U; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4747.211 Safari/537.36
      
# Setting the static fields
fields:
    # Setting the username that we want to brute-force
    - name: username  # the input name
      value: admin
      
# Adding error message if the password is incorrect
# p.s. this will be ignored if you will add on_pass option
on_fail:
    message: incorrect password
    
# Adding the successful message, if, for example, we are in the admin panel
on_pass:
    message: Welcome, 
    
# And the last cherry of this tool is crawl option
# This option can help you find the token if there is any and will add it to your request
crawl:
    url: <token-url> # if the token is not located in the original request, then we will set a new one to get the token
    name: token # the name of the field where the token will be passed to the request
    search: "token = '([a-z0-9]{32})" # to find the token, use regex
```
<p>And that's it, now you are a professional cracker.</p>
