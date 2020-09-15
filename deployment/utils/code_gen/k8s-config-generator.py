import yaml

def get_all_values(root, nested_dictionary, result_dict):
    for key, value in nested_dictionary.items():
        if isinstance(value, dict):
            if(root == ''):
                get_all_values(str.upper(key),value,result_dict)
            else:
                get_all_values(str.upper(root+'_'+key), value,result_dict)  
        else:
            formatted_value=value
            if isinstance(value, list):
                formatted_value = ', '.join(value)
            result_dict[str.upper(root+'_'+key)] = formatted_value
    return result_dict

with open('config.yml') as file:
    documents = yaml.full_load(file)
    result = get_all_values('',documents, {})
    
with open('config-gen.yaml', 'w') as config:
    tmp = []
    for k,v in result.items():
        d = {}
        d[k]=v
        tmp.append(d)
    result = {'config': tmp}
    yaml.dump(result, config)