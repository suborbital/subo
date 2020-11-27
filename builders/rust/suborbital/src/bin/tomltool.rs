use std::fs;
use toml;
use simple_error::SimpleError;
use cargo_toml::{Manifest, Product};


fn main() {
	let file = match read_file("Cargo.toml") {
		None => {
			println!("file not found");
			std::process::exit(1);
		},
		Some(contents) => contents,
	};

	let mut cargo = match deserialize(file) {
		Err(err) => {
			println!("failed to deserialize: {}", err);
			std::process::exit(1);
		}
		Ok(cargo) => cargo,
	};

	add_deps(&mut cargo);

	let serialized = match serialize(&cargo) {
		Err(err) => {
			println!("failed to serialize: {}", err);
			std::process::exit(1);
		},
		Ok(val) => val,
	};

	match write_file("./Cargo.toml", &serialized) {
		Err(err) => {
			println!("failed to write_file: {}", err);
			std::process::exit(1);
		},
		Ok(_) => {},
	};

	println!("{}", cargo.package.unwrap().name)
}

fn add_deps(cargo: &mut Manifest) {
	// let wasm: Dependency = Dependency::Simple(String::from("0.2"));
	// cargo.dependencies.insert(String::from("wasm-bindgen"), wasm);

	let mut lib: Product = cargo.lib.clone().unwrap_or_default();
	lib.crate_type.push(String::from("cdylib"));
	lib.crate_type.push(String::from("rlib"));
	lib.edition = Some(cargo_toml::Edition::E2018);

	cargo.lib = Some(lib);
}

fn read_file(name: &str) -> Option<String> {
	fs::read_to_string(format!("./{}", name)).ok()
}

fn write_file(path: &str, cargo: &String) -> Result<(), SimpleError> {
	fs::write(path, cargo.as_bytes())
		.map_err(|err| SimpleError::from(err))?;

	Ok(())
}

fn deserialize(contents: String) -> Result<Manifest, SimpleError> {
	let cargo: Manifest = cargo_toml::Manifest::from_str(contents.as_str())
		.map_err(|err| SimpleError::from(err))?;

	Ok(cargo)
}

fn serialize(cargo: &Manifest) -> Result<String, SimpleError> {
	let tomlval = toml::Value::try_from(cargo)
		.map_err(|err| SimpleError::from(err))?;

	let serialized = toml::to_string_pretty(&tomlval)
		.map_err(|err| SimpleError::from(err))?;

	Ok(serialized)
}