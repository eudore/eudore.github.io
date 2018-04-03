local file=io.open("data-common.txt","r");
local out=io.open("data-cj.txt","a+")
for line in file:lines() do
	local l=string.match(line,'native (.*) takes')
	if l~= nil then
		out:write(l)
		out:write("\n")
	end
end
out:close()
file:close()
-- cat data-cj.txt | sort | uniq 
